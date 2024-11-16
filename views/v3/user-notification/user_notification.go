package user_notification

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	DiscordEnabled *bool `json:"discord" example:"true" validate:"required"`
	EmailEnabled   *bool `json:"email" example:"true" validate:"required"`
	Events         *bool `json:"events" example:"true" validate:"required"`
	Training       *bool `json:"training" example:"true" validate:"required"`
	Feedback       *bool `json:"feedback" example:"true" validate:"required"`
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
	DiscordEnabled bool `json:"discord" example:"true"`
	EmailEnabled   bool `json:"email" example:"true"`
	Events         bool `json:"events" example:"true"`
	Training       bool `json:"training" example:"true"`
	Feedback       bool `json:"feedback" example:"true"`
}

func NewUserNotificationResponse(r *models.UserNotification) *Response {
	return &Response{
		DiscordEnabled: r.DiscordEnabled,
		EmailEnabled:   r.EmailEnabled,
		Events:         r.Events,
		Training:       r.Training,
		Feedback:       r.Feedback,
	}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetNotificationSettings godoc
// @Summary Get a user notification settings
// @Description Get a user notification settings
// @Tags notification-settings
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/notification-settings [get]
func GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, NewUserNotificationResponse(utils.GetUserNotificationCtx(r)))
}

// UpdateNotificationSettings godoc
// @Summary Update a user's notifications settings
// @Description Update a user notification settings
// @Tags notification-settings
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param request body Request true "Request"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/notification-settings [put]
func UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	userNotification := utils.GetUserNotificationCtx(r)
	userNotification.DiscordEnabled = *req.DiscordEnabled
	userNotification.EmailEnabled = *req.EmailEnabled
	userNotification.Events = *req.Events
	userNotification.Training = *req.Training
	userNotification.Feedback = *req.Feedback

	if err := userNotification.Update(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	utils.Render(w, r, NewUserNotificationResponse(userNotification))
}
