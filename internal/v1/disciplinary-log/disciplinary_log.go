package disciplinary_log

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Entry      string `json:"entry" example:"Changed Preferred OIs to RP" validate:"required"`
	VATUSAOnly bool   `json:"vatusa_only" example:"true"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.DisciplinaryLogEntry
}

func NewDisciplinaryLogEntryResponse(dle *models.DisciplinaryLogEntry) *Response {
	return &Response{DisciplinaryLogEntry: dle}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.DisciplinaryLogEntry == nil {
		return errors.New("disciplinary log entry not found")
	}
	return nil
}

func NewDisciplinaryLogEntryListResponse(dle []models.DisciplinaryLogEntry) []render.Renderer {
	list := []render.Renderer{}
	for idx := range dle {
		list = append(list, NewDisciplinaryLogEntryResponse(&dle[idx]))
	}
	return list
}

// CreateDisciplinaryLogEntry godoc
// @Summary Create a new disciplinary log entry
// @Description Create a new disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param disciplinary_log body Request true "Disciplinary Log Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log [post]
func CreateDisciplinaryLogEntry(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := utils.GetUserCtx(r)

	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	dle := &models.DisciplinaryLogEntry{
		CID:       user.CID,
		Entry:     data.Entry,
		CreatedBy: "System",
	}

	if data.VATUSAOnly {
		dle.VATUSAOnly = true
	} else {
		dle.VATUSAOnly = false
	}

	if err := dle.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// ListDisciplinaryLog godoc
// @Summary List all disciplinary log entries
// @Description List all disciplinary log entries
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param vatusa_only query boolean false "VATUSA Only"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log [get]
func ListDisciplinaryLog(w http.ResponseWriter, r *http.Request) {
	vatusaOnly := r.URL.Query().Get("vatusa_only")
	vatUSA := false
	if vatusaOnly == "true" {
		vatUSA = true
	}

	user := utils.GetUserCtx(r)
	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	dle, err := models.GetAllDisciplinaryLogEntriesByCID(user.CID, vatUSA)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewDisciplinaryLogEntryListResponse(dle)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateDisciplinaryLog godoc
// @Summary Update a disciplinary log entry
// @Description Update a disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param id path string true "Disciplinary Log Entry ID"
// @Param disciplinary_log body Request true "Disciplinary Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [put]
func UpdateDisciplinaryLog(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	dle.Entry = data.Entry

	if data.VATUSAOnly {
		dle.VATUSAOnly = true
	}

	if err := dle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// PatchDisciplinaryLog godoc
// @Summary Patch a disciplinary log entry
// @Description Patch a disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param id path string true "Disciplinary Log Entry ID"
// @Param disciplinary_log body Request true "Disciplinary Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [patch]
func PatchDisciplinaryLog(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Entry != "" {
		dle.Entry = data.Entry
	}

	if data.VATUSAOnly {
		dle.VATUSAOnly = true
	}

	if err := dle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// DeleteDisciplinaryLog godoc
// @Summary Delete a disciplinary log entry
// @Description Delete a disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param id path string true "Disciplinary Log Entry ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [delete]
func DeleteDisciplinaryLog(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)
	if err := dle.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
