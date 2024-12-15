package user_role

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Request struct {
	RoleID     constants.RoleID     `json:"role_id" example:"ATM" validate:"required"`
	FacilityID constants.FacilityID `json:"facility_id" example:"ZDV" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type Response struct {
	Role       constants.RoleID     `json:"role" example:"ATM"`
	FacilityID constants.FacilityID `json:"facility_id" example:"ZDV"`
	CreatedAt  time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

func NewUserRoleResponse(roleID constants.RoleID, facilityId constants.FacilityID, createdAt time.Time) *Response {
	resp := &Response{
		Role:       roleID,
		FacilityID: facilityId,
		CreatedAt:  createdAt,
	}

	return resp
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Role == "" {
		return errors.New("missing required role")
	}
	if res.FacilityID == "" {
		return errors.New("missing required facility-id")
	}
	if res.CreatedAt.IsZero() {
		return errors.New("missing required created-at")
	}

	return nil
}

func NewUserRoleListResponse(userRoles []models.UserRole) []render.Renderer {
	list := []render.Renderer{}
	for idx := range userRoles {
		list = append(list, NewUserRoleResponse(userRoles[idx].RoleID, userRoles[idx].FacilityID, userRoles[idx].CreatedAt))
	}
	return list
}

// GetSelfRoles godoc
// @Summary Get your roles
// @Description Get roles for the user logged in
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roles [get]
func GetSelfRoles(w http.ResponseWriter, r *http.Request) {
	user := utils.GetXUser(r)

	rosters, err := models.GetRostersByCID(user.CID)
	if err != nil {
		log.WithError(err).Errorf("Error getting rosters for user %d", user.CID)
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	roles := []models.UserRole{}
	for _, roster := range rosters {
		roles = append(roles, roster.Roles...)
	}

	if err := render.RenderList(w, r, NewUserRoleListResponse(roles)); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
	}
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
// @Router /user/{cid}/roles [post]
func CreateUserRoles(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	user := utils.GetUserCtx(r)

	roster, err := models.GetRosterByFacilityAndCID(req.FacilityID, user.CID)
	if err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	userRole := &models.UserRole{
		CID:        user.CID,
		RoleID:     req.RoleID,
		FacilityID: req.FacilityID,
		RosterID:   roster.ID,
	}

	if err := userRole.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewUserRoleResponse(userRole.RoleID, userRole.FacilityID, userRole.CreatedAt))

	// Create notification
	notification := &models.Notification{
		CID:      user.CID,
		Category: "Administration",
		Title:    "Role Added",
		Body:     "You have been added to the " + string(userRole.RoleID) + " role at " + string(userRole.FacilityID),
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := notification.Create(); err != nil {
		log.WithError(err).Errorf("Error creating notification for user %d", user.CID)
	}
}

// DeleteUserRoles godoc
// @Summary Remove a user role
// @Description Remove a user role
// @Tags user-roles
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param role_id path string true "Role ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roles/{role_id} [delete]
func DeleteUserRoles(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)
	role := utils.GetUserRoleCtx(r)

	if user.CID != role.CID {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := role.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Create notification
	notification := &models.Notification{
		CID:      user.CID,
		Category: "Administration",
		Title:    "Role Removed",
		Body:     "You have been removed from the " + string(role.RoleID) + " role at " + string(role.FacilityID),
		ExpireAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := notification.Create(); err != nil {
		log.WithError(err).Errorf("Error creating notification for user %d", user.CID)
	}
}
