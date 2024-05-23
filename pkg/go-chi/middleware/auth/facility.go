package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanEditFacility(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		fac := utils.GetFacilityCtx(r)

		if !utils.CanEditFacility(requestingUser, fac) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
