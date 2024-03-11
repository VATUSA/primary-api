package utils

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// These functions grab the {cid} or {id} in the given route and return the associated object
func GetUserCtx(r *http.Request) *models.User {
	return r.Context().Value("user").(*models.User)
}

func GetUserFlagCtx(r *http.Request) *models.UserFlag {
	return r.Context().Value("userFlag").(*models.UserFlag)
}

func GetRosterCtx(r *http.Request) *models.Roster {
	return r.Context().Value("roster").(*models.Roster)
}

func GetFacilityCtx(r *http.Request) (*models.Facility, error) {
	id := chi.URLParam(r, "FacilityID")

	fac := &models.Facility{
		ID: id,
	}

	if err := fac.Get(); err != nil {
		return nil, err
	}

	return fac, nil
}

func GetUserRoleCtx(r *http.Request) *models.UserRole {
	return r.Context().Value("userRole").(*models.UserRole)
}

func GetRosterRequestCtx(r *http.Request) *models.RosterRequest {
	return r.Context().Value("rosterRequest").(*models.RosterRequest)
}

func GetRatingChangeCtx(r *http.Request) *models.RatingChange {
	return r.Context().Value("ratingChange").(*models.RatingChange)
}

func GetActionLogCtx(r *http.Request) *models.ActionLogEntry {
	return r.Context().Value("actionLog").(*models.ActionLogEntry)
}

func GetDisciplinaryLogCtx(r *http.Request) *models.DisciplinaryLogEntry {
	return r.Context().Value("disciplinaryLog").(*models.DisciplinaryLogEntry)
}
