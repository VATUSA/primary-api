package middleware

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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

// FIXME - This is a copy of the request struct in the user_role.go file. This should be moved to a shared location and imported.

type UserRoleRequest struct {
	RoleID     constants.RoleID     `json:"role_id" example:"ATM" validate:"required"`
	FacilityID constants.FacilityID `json:"facility_id" example:"ZDV" validate:"required"`
}

func (req *UserRoleRequest) Validate() error {
	return validator.New().Struct(req)
}

func (req *UserRoleRequest) Bind(r *http.Request) error {
	return nil
}

func CanAddRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := GetSelfUser(r)
		user := utils.GetUserCtx(r)

		if !models.IsValidUser(user.CID) {
			utils.Render(w, r, utils.ErrInvalidCID)
			return
		}

		req := UserRoleRequest{}
		if err := render.Bind(r, &req); err != nil {
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
