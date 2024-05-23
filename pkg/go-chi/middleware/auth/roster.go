package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanEditRoster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		fac := utils.GetFacilityCtx(r)

		if !utils.IsVATUSAStaff(requestingUser) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		if !utils.CanEditFacility(requestingUser, fac) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}
	})
}
