package event

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// CreateEvent godoc
// @Summary Create an Event
// @Description Create an Event
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param event body EventRequest true "Event"
// @Success 201 {object} EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events [post]
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	req := &EventRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	fac := utils.GetFacilityCtx(r)
	if !slices.Contains(req.Facilities, fac.ID) {
		log.Errorf("Facility %s not in facilities list", fac.ID)
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	ev := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		BannerURL:   req.BannerURL,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Fields:      req.Fields,
		Facilities:  req.Facilities,
	}

	if err := ev.Create(); err != nil {
		log.WithError(err).Error("Error creating event")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusCreated)
	utils.Render(w, r, NewEventResponse(ev))
}

// GetAllEvents godoc
// @Summary Get All Events
// @Description Get All Events (Paginated, default 10, limit 25)
// @Tags event
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} []EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /events [get]
func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	if page == 0 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	if limit < 1 || limit > 25 {
		limit = 10
	}

	events, err := models.GetEventsFiltered(int(page), int(limit), "", time.Now())
	if err != nil {
		log.WithError(err).Error("Error getting events")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventListResponse(events)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// GetEvents godoc
// @Summary Get Events
// @Description Get Events by Facility
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} []EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events [get]
func GetEvents(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	events, err := models.GetEventsFiltered(0, 1000, fac.ID, time.Now())
	if err != nil {
		log.WithError(err).Error("Error getting events")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewEventListResponse(events)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// GetEvent godoc
// @Summary Get Event
// @Description Get Event by ID
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID} [get]
func GetEvent(w http.ResponseWriter, r *http.Request) {
	ev := utils.GetEventCtx(r)
	utils.Render(w, r, NewEventResponse(ev))
}

// UpdateEvent godoc
// @Summary Update Event
// @Description Update Event by ID
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param event body EventRequest true "Event"
// @Success 200 {object} EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID} [put]
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ev := utils.GetEventCtx(r)

	req := &EventRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	ev.Title = req.Title
	ev.Description = req.Description
	ev.BannerURL = req.BannerURL
	ev.StartDate = req.StartDate
	ev.EndDate = req.EndDate
	ev.Fields = req.Fields
	ev.Facilities = req.Facilities

	if err := ev.Update(); err != nil {
		log.WithError(err).Error("Error updating event")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewEventResponse(ev))
}

// PatchEvent godoc
// @Summary Patch Event
// @Description Patch Event by ID
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Param event body EventRequest true "Event"
// @Success 200 {object} EventResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID} [patch]
func PatchEvent(w http.ResponseWriter, r *http.Request) {
	ev := utils.GetEventCtx(r)

	req := &EventRequest{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if req.Title != "" {
		ev.Title = req.Title
	}
	if req.Description != "" {
		ev.Description = req.Description
	}
	if req.BannerURL != "" {
		ev.BannerURL = req.BannerURL
	}
	if !req.StartDate.IsZero() {
		ev.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		ev.EndDate = req.EndDate
	}
	if len(req.Fields) > 0 {
		ev.Fields = req.Fields
	}
	if len(req.Facilities) > 0 {
		ev.Facilities = req.Facilities
	}

	if err := ev.Update(); err != nil {
		log.WithError(err).Error("Error updating event")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewEventResponse(ev))
}

// DeleteEvent godoc
// @Summary Delete Event
// @Description Delete Event by ID
// @Tags event
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param EventID path string true "Event ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/events/{EventID} [delete]
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ev := utils.GetEventCtx(r)

	if err := ev.Delete(); err != nil {
		log.WithError(err).Error("Error deleting event")
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Response(r, http.StatusNoContent)
}
