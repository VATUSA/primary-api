package utils

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"net/http"
)

// These functions grab the {cid} or {id} in the given route and return the associated object
func GetUserCtx(r *http.Request) *models.User {
	return r.Context().Value("user").(*models.User)
}

func GetUserFlagCtx(r *http.Request) *models.UserFlag {
	return r.Context().Value("userFlag").(*models.UserFlag)
}
