package middleware

import (
	user_role "github.com/VATUSA/primary-api/external/v3/user-role"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"net/http"
)

func CanViewUser(next http.Handler) http.Handler {
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

func CanEditUser(next http.Handler) http.Handler {
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

func CanAddRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !models.IsValidUser(user.CID) {
			utils.Render(w, r, utils.ErrInvalidCID)
			return
		}

		req := &user_role.Request{}
		if err := render.Bind(r, req); err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			utils.Render(w, r, utils.ErrBadRequest)
			return
		}

		if !req.RoleID.IsValidRole() {
			utils.Render(w, r, utils.ErrInvalidRole)
			return
		}

		if !utils.CanAddRole(requestingUser, req.RoleID, req.FacilityID) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CanDeleteRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !models.IsValidUser(user.CID) {
			utils.Render(w, r, utils.ErrInvalidCID)
			return
		}

		role := utils.GetUserRoleCtx(r)

		if !utils.CanAddRole(requestingUser, role.RoleID, role.FacilityID) {
			utils.Render(w, r, utils.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
