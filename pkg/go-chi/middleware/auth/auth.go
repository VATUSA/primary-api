package middleware

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/cookie"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
	"strconv"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("VATUSA")
		if err != nil {
			// Set x-guest context to true
			ctx := context.WithValue(r.Context(), utils.XGuest{}, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		auth := make(map[string]string)
		if err := cookie.CookieStore.Decode("VATUSA", authCookie.Value, &auth); err != nil {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		cid, err := strconv.ParseUint(auth["cid"], 10, 64)
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		user := &models.User{CID: uint(cid)}
		if err := user.Get(); err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.XUser{}, user)
		ctx = context.WithValue(ctx, utils.XGuest{}, false)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func NotGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guest := utils.GetXGuest(r)
		if guest {
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

// TODO - Implement properly
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

func HasRoleInFacility(w http.ResponseWriter, r *http.Request, facility constants.FacilityID, roles ...constants.RoleID) bool {
	user := GetSelfUser(r)

	for _, role := range roles {
		if models.HasRoleAtFacility(user, role, facility) {
			return true
		}
	}

	return false
}

func GetSelfUser(r *http.Request) *models.User {
	user := utils.GetXUser(r)

	return user
}
