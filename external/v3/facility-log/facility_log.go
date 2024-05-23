package facility_log

import (
	"errors"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Facility constants.FacilityID `json:"facility" example:"ZDV" validate:"required,len=3"`
	Entry    string               `json:"entry" example:"Changed Preferred OIs to RP" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.FacilityLogEntry
}

func NewFacilityLogEntryResponse(fle *models.FacilityLogEntry) *Response {
	return &Response{FacilityLogEntry: fle}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.FacilityLogEntry == nil {
		return errors.New("facility log entry not found")
	}
	return nil
}

func NewFacilityLogEntryListResponse(fle []models.FacilityLogEntry) []render.Renderer {
	list := []render.Renderer{}
	for _, f := range fle {
		list = append(list, NewFacilityLogEntryResponse(&f))
	}
	return list
}

// CreateFacilityLogEntry godoc
// @Summary Create a new facility log entry
// @Description Create a new facility log entry
// @Tags facility-log
// @Accept  json
// @Produce  json
// @Param facility_log body Request true "Facility Log Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility-log [post]
func CreateFacilityLogEntry(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	requestingUser := middleware.GetSelfUser(r)

	fle := &models.FacilityLogEntry{
		Facility:  data.Facility,
		Entry:     data.Entry,
		CreatedBy: fmt.Sprintf("%d", requestingUser.CID),
	}

	if err := fle.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewFacilityLogEntryResponse(fle))
}

// ListFacilityLog godoc
// @Summary List facility log entries
// @Description List facility log entries
// @Tags facility-log
// @Accept  json
// @Produce  json
// @Param facility query string false "Facility ID"
// @Success 200 {object} Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility-log [get]
func ListFacilityLog(w http.ResponseWriter, r *http.Request) {
	facId := r.URL.Query().Get("facility")
	if facId != "" {
		if !models.IsValidFacility(constants.FacilityID(facId)) {
			utils.Render(w, r, utils.ErrInvalidFacility)
			return
		}
		fle, err := models.GetAllFacilityLogEntriesByFacility(constants.FacilityID(facId))
		if err != nil {
			utils.Render(w, r, utils.ErrInternalServer)
			return
		}

		if err := render.RenderList(w, r, NewFacilityLogEntryListResponse(fle)); err != nil {
			utils.Render(w, r, utils.ErrRender(err))
			return
		}
		return
	}

	fle, err := models.GetAllFacilityLogEntries()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewFacilityLogEntryListResponse(fle)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateFacilityLog godoc
// @Summary Update a facility log entry
// @Description Update a facility log entry
// @Tags facility-log
// @Accept  json
// @Produce  json
// @Param id path string true "Facility Log Entry ID"
// @Param facility_log body Request true "Facility Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility-log/{id} [put]
func UpdateFacilityLog(w http.ResponseWriter, r *http.Request) {
	fle := utils.GetFacilityLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	fle.Facility = data.Facility
	fle.Entry = data.Entry

	if err := fle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFacilityLogEntryResponse(fle))
}

// PatchFacilityLog godoc
// @Summary Patch a facility log entry
// @Description Patch a facility log entry
// @Tags facility-log
// @Accept  json
// @Produce  json
// @Param id path string true "Facility Log Entry ID"
// @Param facility_log body Request true "Facility Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility-log/{id} [patch]
func PatchFacilityLog(w http.ResponseWriter, r *http.Request) {
	fle := utils.GetFacilityLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Facility != "" {
		if !models.IsValidFacility(data.Facility) {
			utils.Render(w, r, utils.ErrInvalidFacility)
			return
		}
		fle.Facility = data.Facility
	}
	if data.Entry != "" {
		fle.Entry = data.Entry

	}

	if err := fle.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFacilityLogEntryResponse(fle))
}

// DeleteFacilityLog godoc
// @Summary Delete a facility log entry
// @Description Delete a facility log entry
// @Tags facility-log
// @Accept  json
// @Produce  json
// @Param id path string true "Facility Log Entry ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /facility-log/{id} [delete]
func DeleteFacilityLog(w http.ResponseWriter, r *http.Request) {
	fle := utils.GetFacilityLogCtx(r)

	if err := fle.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
