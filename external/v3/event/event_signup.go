package event

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// CreateEventSignup godoc
// @Summary Create an Event Signup
// @Description Create an Event Signup
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param event body EventSignupRequest true "Event Signup"
// @Success 201 {object} EventSignupResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/signups [post]
func CreateEventSignup(w http.ResponseWriter, r *http.Request) {
	req := &EventSignupRequest{}
	if err := req.Bind(r); err != nil {
		log.WithError(err).Error("Error binding request")
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.WithError(err).Error("Error validating request")
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	event := utils.GetEventCtx(r)
	signup := &models.EventSignup{
		EventID:    event.ID,
		PositionID: req.PositionID,
		CID:        req.CID,
		Shift:      req.Shift,
	}

	if err := signup.Create(); err != nil {
		log.WithError(err).Error("Error creating event signup")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusCreated)
	utils.Render(w, r, NewEventSignupResponse(signup))
}

// GetEventSignup godoc
// @Summary Get an Event Signup
// @Description Get an Event Signup
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventSignupID path string true "Event Signup ID"
// @Success 200 {object} EventSignupResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/signups/{EventSignupID} [get]
func GetEventSignup(w http.ResponseWriter, r *http.Request) {
	signup := utils.GetEventSignupCtx(r)
	utils.Render(w, r, NewEventSignupResponse(signup))
}

// GetEventSignups godoc
// @Summary Get Event Signups
// @Description Get Event Signups
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Success 200 {object} []EventSignupResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/signups [get]
func GetEventSignups(w http.ResponseWriter, r *http.Request) {
	event := utils.GetEventCtx(r)
	signups, err := models.GetEventSignupFiltered(event.ID)
	if err != nil {
		log.WithError(err).Error("Error getting event signups")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventSignupListResponse(signups)); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
}

// DeleteEventSignup godoc
// @Summary Delete an Event Signup
// @Description Delete an Event Signup
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param EventSignupID path string true "Event Signup ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID}/signups/{EventSignupID} [delete]
func DeleteEventSignup(w http.ResponseWriter, r *http.Request) {
	signup := utils.GetEventSignupCtx(r)

	if err := signup.Delete(); err != nil {
		log.WithError(err).Error("Error deleting event signup")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusNoContent)
}
