package middleware

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

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
		credentials := GetCredentials(r)
		targetUser := utils.GetUserCtx(r)
		if targetUser.Flags.NoStaffRole {
			utils.Render(w, r, utils.ErrInvalidRequest(errors.New("user is not allowed to hold any staff roles")))
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

		if credentials.User != nil {
			if !utils.CanAddRole(credentials.User, req.RoleID, req.FacilityID) {
				utils.Render(w, r, utils.ErrForbidden)
				return
			}
		}

		if credentials.Facility != nil {
			// TODO - possibly allow this in the future
		}

		next.ServeHTTP(w, r)
	})
}

func CanDeleteRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestingUser := utils.GetXUser(r)
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
