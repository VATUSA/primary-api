package disciplinary_log

import (
	"fmt"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Entry      string `json:"entry" example:"Misconduct in discord" validate:"required"`
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
		return nil
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
// @Param cid path int true "User CID"
// @Param action_log body Request true "Disciplinary Log Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log [post]
func CreateDisciplinaryLogEntry(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
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

	createdBy := ""
	if self := utils.GetXUser(r); self != nil {
		createdBy = fmt.Sprintf("%d", self.CID)
	} else {
		createdBy = string(utils.GetXFacility(r).ID)
	}

	dle := &models.DisciplinaryLogEntry{
		CID:        user.CID,
		Entry:      data.Entry,
		VATUSAOnly: data.VATUSAOnly,
		CreatedBy:  createdBy,
	}

	if err := dle.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// GetDisciplinaryLog godoc
// @Summary Get all disciplinary log entries
// @Description List all disciplinary log entries
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param vatusa_only query bool false "VATUSA Only"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log [get]
func GetDisciplinaryLog(w http.ResponseWriter, r *http.Request) {
	vatusaOnly := r.URL.Query().Get("vatusa_only")

	user := utils.GetUserCtx(r)
	dle, err := models.GetAllDisciplinaryLogEntriesByCID(user.CID, vatusaOnly == "true")
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewDisciplinaryLogEntryListResponse(dle)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateDisciplinaryLogEntry godoc
// @Summary Update a disciplinary log entry
// @Description Update a disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param id path int true "Disciplinary Log Entry ID"
// @Param disciplinary_log body Request true "Disciplinary Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [put]
func UpdateDisciplinaryLogEntry(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	dle.Entry = data.Entry
	dle.VATUSAOnly = data.VATUSAOnly

	if self := utils.GetXUser(r); self != nil {
		dle.UpdatedBy = fmt.Sprintf("%d", self.CID)
	} else {
		dle.UpdatedBy = string(utils.GetXFacility(r).ID)
	}

	if err := dle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// PatchDisciplinaryLogEntry godoc
// @Summary Patch an disciplinary log entry
// @Description Patch an disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Disciplinary Log Entry ID"
// @Param disciplinary_log body Request true "Disciplinary Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [patch]
func PatchDisciplinaryLogEntry(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Entry != "" {
		dle.Entry = data.Entry
	}
	if data.VATUSAOnly {
		dle.VATUSAOnly = data.VATUSAOnly
	}

	if self := utils.GetXUser(r); self != nil {
		dle.UpdatedBy = fmt.Sprintf("%d", self.CID)
	} else {
		dle.UpdatedBy = string(utils.GetXFacility(r).ID)
	}

	if err := dle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDisciplinaryLogEntryResponse(dle))
}

// DeleteDisciplinaryLogEntry godoc
// @Summary Delete a disciplinary log entry
// @Description Delete a disciplinary log entry
// @Tags disciplinary-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Disciplinary Log Entry ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/disciplinary-log/{id} [delete]
func DeleteDisciplinaryLogEntry(w http.ResponseWriter, r *http.Request) {
	dle := utils.GetDisciplinaryLogCtx(r)

	if err := dle.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
