package middleware

import (
	"context"
	"errors"
	"github.com/VATUSA/primary-api/pkg/cookie"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
	"strconv"
)

func HasCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("VATUSA")
		if errors.Is(err, http.ErrNoCookie) {
			// Set x-guest context to true
			ctx := context.WithValue(r.Context(), utils.XGuest{}, true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
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

func HasAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			//ctx := context.WithValue(r.Context(), utils.XGuest{}, true)
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		facility, err := models.GetFacilityByAPIKey(apiKey)
		if err != nil {
			//ctx := context.WithValue(r.Context(), utils.XGuest{}, true)
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		ctx := context.WithValue(r.Context(), utils.XFacility{}, facility)
		ctx = context.WithValue(ctx, utils.XGuest{}, false)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NotGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guest := utils.GetXGuest(r)
		if guest {
			utils.Render(w, r, utils.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type Credentials struct {
	User     *models.User
	Facility *models.Facility
}

func GetCredentials(r *http.Request) *Credentials {
	user := utils.GetXUser(r)
	facility := utils.GetXFacility(r)

	return &Credentials{
		User:     user,
		Facility: facility,
	}
}
