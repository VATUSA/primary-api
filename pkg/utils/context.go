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
