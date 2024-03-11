package user

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Request struct {
	CID              uint      `json:"cid" example:"1293257" validate:"required"`
	FirstName        string    `json:"first_name" example:"Raaj" validate:"required"`
	LastName         string    `json:"last_name" example:"Patel" validate:"required"`
	PreferredName    string    `json:"preferred_name" example:"Raaj"`
	Email            string    `json:"email" example:"vatusa6@vatusa.net" validate:"required,email"`
	PreferredOIs     string    `json:"preferred_ois" example:"RP"`
	PilotRating      uint      `json:"pilot_rating" example:"1" validate:"required"`
	ControllerRating uint      `json:"controller_rating" example:"1" validate:"required"`
	DiscordID        string    `json:"discord_id" example:"1234567890"`
	LastCertSync     time.Time `json:"last_cert_sync" example:"2021-01-01T00:00:00Z"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.User
}

func NewUserResponse(user *models.User) *Response {
	resp := &Response{User: user}

	return resp
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.User == nil {
		return errors.New("missing required user")
	}
	return nil
}

func NewUserListResponse(users []models.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(&user))
	}
	return list
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body Request true "User"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := &models.User{
		CID:              req.CID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		PreferredName:    req.PreferredName,
		Email:            req.Email,
		PreferredOIs:     req.PreferredOIs,
		PilotRating:      req.PilotRating,
		ControllerRating: req.ControllerRating,
		DiscordID:        req.DiscordID,
	}
	if err := user.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	userflag := &models.UserFlag{
		CID:            user.CID,
		NoStaffRole:    false,
		NoVisiting:     false,
		NoTransferring: false,
		NoTraining:     false,
	}
	if err := userflag.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return

	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewUserResponse(user))
}

// GetUser godoc
// @Summary Get a user
// @Description Get a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	utils.Render(w, r, NewUserResponse(user))
}

// ListUsers godoc
// @Summary List users
// @Description List users
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewUserListResponse(users)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param user body Request true "User"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user.CID = req.CID
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PreferredName = req.PreferredName
	user.Email = req.Email
	user.PreferredOIs = req.PreferredOIs
	user.PilotRating = req.PilotRating
	user.ControllerRating = req.ControllerRating
	user.DiscordID = req.DiscordID

	if err := user.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewUserResponse(user))
}

// PatchUser godoc
// @Summary Patch a user
// @Description Patch a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param user body Request true "User"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [patch]
func PatchUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PreferredName != "" {
		user.PreferredName = req.PreferredName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PreferredOIs != "" {
		user.PreferredOIs = req.PreferredOIs
	}
	if req.DiscordID != "" {
		user.DiscordID = req.DiscordID
	}

	if err := user.Update(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Render(w, r, NewUserResponse(user))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user
// @Tags user
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	if err := user.Delete(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
