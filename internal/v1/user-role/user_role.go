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

func NewUserRoleListResponse(userRoles []models.UserRole) []render.Renderer {
	list := []render.Renderer{}
	for idx := range userRoles {
		list = append(list, NewUserRoleResponse(&userRoles[idx]))
	}
	return list
}

// CreateUserRoles godoc
// @Summary Create a new user role
// @Description Create a new user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param user_role body Request true "User Role"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/role [post]
func CreateUserRoles(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(*r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := utils.GetUserCtx(r)

	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !req.RoleID.IsValidRole() {
		utils.Render(w, r, utils.ErrInvalidRole)
		return
	}

	roster, err := models.GetRosterByFacilityAndCID(req.FacilityID, user.CID)
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	userRole := &models.UserRole{
		CID:        user.CID,
		RoleID:     req.RoleID,
		FacilityID: req.FacilityID,
		RosterID:   roster.ID,
	}

	if err := userRole.Create(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewUserRoleResponse(userRole))
}

// GetUserRole godoc
// @Summary Get a user role
// @Description Get a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "User Role ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/role/{id} [get]
func GetUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := utils.GetUserRoleCtx(r)

	utils.Render(w, r, NewUserRoleResponse(userRole))
}

// ListUserRoles godoc
// @Summary List user roles
// @Description List user roles
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/role [get]
func ListUserRoles(w http.ResponseWriter, r *http.Request) {
	userRoles, err := models.GetAllUserRoles()
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewUserRoleListResponse(userRoles)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// DeleteUserRole godoc
// @Summary Delete a user role
// @Description Delete a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Success 204
// @Param cid path int true "User CID"
// @Param id path int true "User Role ID"
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/role/{id} [delete]
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	userRole := utils.GetUserRoleCtx(r)

	if err := userRole.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
