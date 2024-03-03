package middleware

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"net/http"
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

func HasRoles(roles ...constants.RoleID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetSelfUser(r)

			for _, role := range roles {
				for _, userRole := range user.Roles {
					if role == userRole.RoleID {
						next.ServeHTTP(w, r)
						return
					}
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

func HasRoleInFacility(w http.ResponseWriter, r *http.Request, facility constants.FacilityID, role ...constants.RoleID) bool {
	user := GetSelfUser(r)

	for _, userRole := range user.Roles {
		if userRole.FacilityID == facility {
			for _, r := range role {
				if userRole.RoleID == r {
					return true
				}
			}
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
