package middleware

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanViewFacilityLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)

		facParam := r.URL.Query().Get("facility")

		if facParam == "" {
			if !utils.IsVATUSAStaff(requestingUser) {
				utils.Render(w, r, utils.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		fac := &models.Facility{ID: constants.FacilityID(facParam)}
		err := fac.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		if !utils.CanEditFacility(requestingUser, fac) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CanEditFacilityLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		facLog := utils.GetFacilityLogCtx(r)

		fac := &models.Facility{ID: facLog.Facility}
		err := fac.Get()
		if err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		if !utils.CanEditFacility(requestingUser, fac) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
