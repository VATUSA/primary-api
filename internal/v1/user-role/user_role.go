package user_role

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	CID        uint                 `json:"cid" example:"1293257" validate:"required"`
	RoleID     constants.RoleID     `json:"role_id" example:"ATM" validate:"required"`
	FacilityID constants.FacilityID `json:"facility_id" example:"ZDV" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r http.Request) error {
	return nil
}

type Response struct {
	*models.UserRole
}

func NewUserRoleResponse(r *models.UserRole) *Response {
	return &Response{UserRole: r}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.UserRole == nil {
		return errors.New("user role not found")
	}

	return nil
}

func NewUserRoleListResponse(r []models.UserRole) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range r {
		list = append(list, NewUserRoleResponse(&d))
	}
	return list
}

// CreateUserRoles godoc
// @Summary Create a new user role
// @Description Create a new user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param user_role body Request true "User Role"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles [post]
func CreateUserRoles(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(*r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(req.CID) {
		render.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !req.RoleID.IsValidRole() {
		render.Render(w, r, utils.ErrInvalidRole)
		return
	}

	userRole := &models.UserRole{
		CID:        req.CID,
		RoleID:     req.RoleID,
		FacilityID: req.FacilityID,
	}

	if err := userRole.Create(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserRoleResponse(userRole))
}

// GetUserRole godoc
// @Summary Get a user role
// @Description Get a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param id path int true "User Role ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles/{id} [get]
func GetUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := GetUserRoleCtx(r)

	render.Render(w, r, NewUserRoleResponse(userRole))
}

// ListUserRoles godoc
// @Summary List user roles
// @Description List user roles
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles [get]
func ListUserRoles(w http.ResponseWriter, r *http.Request) {
	userRoles, err := models.GetAllUserRoles()
	if err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewUserRoleListResponse(userRoles)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateUserRole godoc
// @Summary Update a user role
// @Description Update a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param user_role body Request true "User Role"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles [put]
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := GetUserRoleCtx(r)

	req := &Request{}
	if err := req.Bind(*r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(req.CID) {
		render.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !req.RoleID.IsValidRole() {
		render.Render(w, r, utils.ErrInvalidRole)
		return
	}

	userRole.CID = req.CID
	userRole.RoleID = req.RoleID
	userRole.FacilityID = req.FacilityID

	if err := userRole.Update(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserRoleResponse(userRole))
}

// PatchUserRole godoc
// @Summary Patch a user role
// @Description Patch a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param user_role body Request true "User Role"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles [patch]
func PatchUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := GetUserRoleCtx(r)

	req := &Request{}
	if err := req.Bind(*r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if req.CID != 0 {
		if !models.IsValidUser(req.CID) {
			render.Render(w, r, utils.ErrInvalidCID)
			return
		}
		userRole.CID = req.CID
	}
	if req.RoleID != "" {
		if !req.RoleID.IsValidRole() {
			render.Render(w, r, utils.ErrInvalidRequest(errors.New("invalid role")))
			return
		}
		userRole.RoleID = req.RoleID
	}
	if req.FacilityID != "" {
		userRole.FacilityID = req.FacilityID
	}

	if err := userRole.Update(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserRoleResponse(userRole))
}

// DeleteUserRole godoc
// @Summary Delete a user role
// @Description Delete a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user-roles [delete]
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := GetUserRoleCtx(r)

	if err := userRole.Delete(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
