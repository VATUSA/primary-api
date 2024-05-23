package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanViewNotifications(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !utils.IsVATUSAStaff(requestingUser) && requestingUser.CID != user.CID {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CanEditNotifications(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)

		if !utils.IsVATUSAStaff(requestingUser) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
