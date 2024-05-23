package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
	"strconv"
)

func CanViewFeedback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		fac := utils.GetFacilityCtx(r)

		cidQuery := r.URL.Query().Get("cid")
		cidInt, err := strconv.Atoi(cidQuery)
		if err != nil {
			utils.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}

		statusQuery := r.URL.Query().Get("status")
		if statusQuery == "" || statusQuery == "approved" {
			if uint(cidInt) != requestingUser.CID {
				utils.Render(w, r, utils.ErrForbidden)
				return
			}
		}

		if !utils.IsVATUSAStaff(requestingUser) && !utils.IsFacilitySeniorStaff(requestingUser, fac.ID) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CanEditFeedback(next http.Handler) http.Handler {
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
