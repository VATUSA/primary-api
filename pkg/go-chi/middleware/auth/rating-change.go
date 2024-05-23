package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanViewRatingChange(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !utils.CanViewUser(requestingUser, user) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CanEditRatingChange(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !utils.IsVATUSAStaff(requestingUser) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		for _, roster := range user.Roster {
			if utils.IsInstructor(requestingUser, roster.Facility) {
				next.ServeHTTP(w, r)
				return
			}
		}

		utils.Render(w, r, utils.ErrForbidden)
	})
}
