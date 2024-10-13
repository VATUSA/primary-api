package event

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// CreateEventPosition godoc
// @Summary Create an Event Position
// @Description Create an Event Position
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param event body EventPositionRequest true "Event Position"
// @Success 201 {object} EventPositionResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/positions [post]
func CreateEventPosition(w http.ResponseWriter, r *http.Request) {
	req := &EventPositionRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	fac := utils.GetFacilityCtx(r)
	position := &models.EventPosition{
		EventID:           utils.GetEventCtx(r).ID,
		Facility:          fac.ID,
		Position:          req.Position,
		Shifts:            false,
		Assignee:          0,
		SecondaryAssignee: 0,
	}

	if req.Shifts != nil {
		position.Shifts = *req.Shifts
	}
	if req.Assignee != nil {
		position.Assignee = *req.Assignee
	}
	if req.SecondaryAssignee != nil {
		position.SecondaryAssignee = *req.SecondaryAssignee
	}

	if err := position.Create(); err != nil {
		log.WithError(err).Error("Error creating event position")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusCreated)
	utils.Render(w, r, NewEventPositionResponse(position))
}

// GetEventPositions godoc
// @Summary Get Event Positions
// @Description Get Event Positions
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Success 200 {object} []EventPositionResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/positions [get]
func GetEventPositions(w http.ResponseWriter, r *http.Request) {
	positions, err := models.GetEventPositionsFiltered(utils.GetEventCtx(r).ID, utils.GetFacilityCtx(r).ID)
	if err != nil {
		log.WithError(err).Error("Error getting event positions")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventPositionListResponse(positions)); err != nil {
		log.WithError(err).Error("Error rendering event positions")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
}

// GetEventPosition godoc
// @Summary Get Event Position
// @Description Get Event Position
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventPositionID path string true "Event Position ID"
// @Success 200 {object} EventPositionResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/positions/{EventPositionID} [get]
func GetEventPosition(w http.ResponseWriter, r *http.Request) {
	position := utils.GetEventPositionCtx(r)
	utils.Render(w, r, NewEventPositionResponse(position))
}

// PatchEventPosition godoc
// @Summary Patch an Event Position
// @Description Patch an Event Position
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventPositionID path string true "Event Position ID"
// @Param event body EventPositionRequest true "Event Position"
// @Success 200 {object} EventPositionResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/positions/{EventPositionID} [patch]
func PatchEventPosition(w http.ResponseWriter, r *http.Request) {
	position := utils.GetEventPositionCtx(r)

	req := &EventPositionRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if req.Position != "" {
		position.Position = req.Position
	}
	if req.Assignee != nil {
		position.Assignee = *req.Assignee
	}
	if req.Shifts != nil {
		position.Shifts = *req.Shifts
	}
	if req.SecondaryAssignee != nil {
		position.SecondaryAssignee = *req.SecondaryAssignee
	}

	if err := position.Update(); err != nil {
		log.WithError(err).Error("Error updating event position")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewEventPositionResponse(position))
}

// DeleteEventPosition godoc
// @Summary Delete an Event Position
// @Description Delete an Event Position
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventPositionID path string true "Event Position ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/positions/{EventPositionID} [delete]
func DeleteEventPosition(w http.ResponseWriter, r *http.Request) {
	position := utils.GetEventPositionCtx(r)
	if err := position.Delete(); err != nil {
		log.WithError(err).Error("Error deleting event position")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusNoContent)
}
