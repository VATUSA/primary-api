package event

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetEventRouting godoc
// @Summary Get Event Routing
// @Description Get Event Routing
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Success 200 {object} []EventRoutingResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/routing [get]
func GetEventRouting(w http.ResponseWriter, r *http.Request) {
	event := utils.GetEventCtx(r)
	routing, err := models.GetEventRoutingFiltered(event.ID)
	if err != nil {
		log.WithError(err).Error("Error getting event routing")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventRoutingListResponse(routing)); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
}

// CreateEventRouting godoc
// @Summary Create an Event Routing
// @Description Create an Event Routing
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param event body EventRoutingRequest true "Event Routing"
// @Success 201 {object} EventRoutingResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/routing [post]
func CreateEventRouting(w http.ResponseWriter, r *http.Request) {
	req := &EventRoutingRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	event := utils.GetEventCtx(r)
	routing := &models.EventRouting{
		EventID:     event.ID,
		Origin:      req.Origin,
		Destination: req.Destination,
		Routing:     req.Routing,
		Notes:       req.Notes,
	}

	if err := routing.Create(); err != nil {
		log.WithError(err).Error("Error creating event routing")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusCreated)
	utils.Render(w, r, NewEventRoutingResponse(routing))
}

// PatchEventRouting godoc
// @Summary Patch an Event Routing
// @Description Patch an Event Routing
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventRoutingID path string true "Event Routing ID"
// @Param event body EventRoutingRequest true "Event Routing"
// @Success 200 {object} EventRoutingResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/routing/{EventRoutingID} [patch]
func PatchEventRouting(w http.ResponseWriter, r *http.Request) {
	routing := utils.GetEventRoutingCtx(r)

	req := &EventRoutingRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if req.Origin != "" {
		routing.Origin = req.Origin
	}
	if req.Destination != "" {
		routing.Destination = req.Destination
	}
	if req.Routing != "" {
		routing.Routing = req.Routing
	}
	if req.Notes != "" {
		routing.Notes = req.Notes
	}

	if err := routing.Update(); err != nil {
		log.WithError(err).Error("Error patching event routing")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewEventRoutingResponse(routing))
}

// DeleteEventRouting godoc
// @Summary Delete an Event Routing
// @Description Delete an Event Routing
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventRoutingID path string true "Event Routing ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/routing/{EventRoutingID} [delete]
func DeleteEventRouting(w http.ResponseWriter, r *http.Request) {
	routing := utils.GetEventRoutingCtx(r)
	if err := routing.Delete(); err != nil {
		log.WithError(err).Error("Error deleting event routing")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusNoContent)
}
