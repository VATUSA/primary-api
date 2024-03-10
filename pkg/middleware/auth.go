package middleware

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"net/http"
	"strconv"
)

func NotGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guest := r.Header.Get("x-guest")
		if guest == "true" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HasRole(roles ...constants.RoleID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetSelfUser(r)

			for _, role := range roles {
				if models.HasRole(user, role) {
					next.ServeHTTP(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HasAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HasRoleInFacility(w http.ResponseWriter, r *http.Request, facility string, roles ...constants.RoleID) bool {
	user := GetSelfUser(r)

	for _, role := range roles {
		if models.HasRoleAtFacility(user, role, facility) {
			return true
		}
	}

	return false
}

func GetSelfUser(r *http.Request) *models.User {
	xuser := r.Header.Get("x-user")

	user := &models.User{}
	if err := json.Unmarshal([]byte(xuser), user); err != nil {
		return nil
	}

	return user
}

func GetSelfUserCID(r *http.Request) string {
	xuser := r.Header.Get("x-user")

	user := &models.User{}
	if err := json.Unmarshal([]byte(xuser), user); err != nil {
		return ""
	}

	return strconv.Itoa(int(user.CID))
}

func IsGuest(r *http.Request) bool {
	return r.Header.Get("x-guest") == "true"
}
