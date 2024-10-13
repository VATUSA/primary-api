package event

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetEventTemplates godoc
// @Summary Get Event Templates
// @Description Get Event Templates by Facility
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} []EventTemplateResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/event-templates [get]
func GetEventTemplates(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	eventTemplates, err := models.GetEventTemplatesFiltered(fac.ID)
	if err != nil {
		log.WithError(err).Error("Error getting event templates")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventTemplateListResponse(eventTemplates)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateEventTemplate godoc
// @Summary Update an Event Template
// @Description Update an Event Template
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventTemplateID path string true "Event Template ID"
// @Param event body EventTemplateRequest true "Event Template"
// @Success 200 {object} EventTemplateResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/event-templates/{EventTemplateID} [put]
func UpdateEventTemplate(w http.ResponseWriter, r *http.Request) {
	et := utils.GetEventTemplateCtx(r)

	req := &EventTemplateRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	et.Title = req.Title
	et.Positions = req.Positions
	et.Facilities = req.Facilities
	et.Fields = req.Fields
	et.Shifts = req.Shifts

	if err := et.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewEventTemplateResponse(et))
}

// DeleteEventTemplate godoc
// @Summary Delete an Event Template
// @Description Delete an Event Template
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventTemplateID path string true "Event Template ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/event-templates/{EventTemplateID} [delete]
func DeleteEventTemplate(w http.ResponseWriter, r *http.Request) {
	et := utils.GetEventTemplateCtx(r)

	if err := et.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusNoContent)
}
