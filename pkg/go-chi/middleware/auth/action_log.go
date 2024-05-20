package middleware

import (
	"github.com/VATUSA/primary-api/pkg/utils"
	"net/http"
)

func CanReadActionLog(next http.Handler) http.Handler {
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

func CanEditActionLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !utils.CanEditUser(requestingUser, user) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
