package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanEditNews(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		fac := utils.GetFacilityCtx(r)

		if !utils.IsVATUSAStaff(requestingUser) && !utils.IsFacilitySeniorStaff(requestingUser, fac.ID) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
