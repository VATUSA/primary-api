package user

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Request struct {
	PreferredName string `json:"preferred_name" example:"Raaj" validate:"required"`
	PreferredOIs  string `json:"preferred_ois" example:"RP" validate:"required"`
	DiscordID     string `json:"discord_id" example:"1234567890" validate:"required"`
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

// TODO - Dont send back entire struct
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
	for idx := range users {
		list = append(list, NewUserResponse(&users[idx]))
	}
	return list
}

// GetSelf godoc
// @Summary Get your user
// @Description Get information for the user logged in
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/ [get]
func GetSelf(w http.ResponseWriter, r *http.Request) {
	user := utils.GetXUser(r)

	utils.Render(w, r, NewUserResponse(user))
}

// GetUser godoc
// @Summary Get a specific user
// @Description Get information for the user
// @Tags user
// @Accept  json
// @Produce  json
// @Param CID path int true "CID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	utils.Render(w, r, NewUserResponse(user))
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update information for the user
// @Tags user
// @Accept  json
// @Produce  json
// @Param CID path int true "CID"
// @Param user body Request true "User"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	user.PreferredName = req.PreferredName
	user.PreferredOIs = strings.ToUpper(req.PreferredOIs)
	user.DiscordID = req.DiscordID

	if err := user.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewUserResponse(user))
}

// PatchUser godoc
// @Summary Patch a user
// @Description Patch information for the user
// @Tags user
// @Accept  json
// @Produce  json
// @Param CID path int true "CID"
// @Param user body Request true "User"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid} [patch]
func PatchUser(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if req.PreferredName != "" {
		if req.PreferredName == "-" {
			req.PreferredName = ""
		}
		user.PreferredName = req.PreferredName
	}
	if req.PreferredOIs != "" {
		user.PreferredOIs = strings.ToUpper(req.PreferredOIs)
	}
	if req.DiscordID != "" {
		user.DiscordID = req.DiscordID
	}

	if err := user.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewUserResponse(user))
}
