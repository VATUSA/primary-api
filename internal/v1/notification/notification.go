package notification

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
	CID      uint   `json:"cid" example:"1293257" validate:"required"`
	Category string `json:"category" example:"Training" validate:"required"`
	Title    string `json:"title" example:"Upcoming Training Session" validate:"required"`
	Body     string `json:"body" example:"You have a training session coming up." validate:"required"`
	ExpireAt string `json:"expire_at" example:"2021-01-01T00:00:00Z" validate:"required"`
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
	*models.Notification
}

func NewNotificationResponse(n *models.Notification) *Response {
	return &Response{Notification: n}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Notification == nil {
		return errors.New("notification not found")
	}
	return nil
}

func NewNotificationListResponse(n []models.Notification) []render.Renderer {
	list := []render.Renderer{}
	for idx := range n {
		list = append(list, NewNotificationResponse(&n[idx]))
	}
	return list
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification
// @Tags notification
// @Accept  json
// @Produce  json
// @Param notification body Request true "Notification"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /notification [post]
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	expireAt, err := http.ParseTime(data.ExpireAt)
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	// Make sure expireAt is in the future
	if expireAt.Before(time.Now()) {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("expire_at must be in the future")))
		return
	}

	if !models.IsValidUser(data.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	n := &models.Notification{
		CID:      data.CID,
		Category: data.Category,
		Title:    data.Title,
		Body:     data.Body,
		ExpireAt: expireAt,
	}

	if err := n.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewNotificationResponse(n))
}

// GetNotification godoc
// @Summary Get a notification
// @Description Get a notification
// @Tags notification
// @Param id path int true "Notification ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Router /notification/{id} [get]
func GetNotification(w http.ResponseWriter, r *http.Request) {
	n := utils.GetNotificationCtx(r)
	utils.Render(w, r, NewNotificationResponse(n))
}

// ListNotifications godoc
// @Summary List all notifications
// @Description List all notifications
// @Tags notification
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /notification [get]
func ListNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := models.GetAllNotifications()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewNotificationListResponse(notifications)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateNotification godoc
// @Summary Update a notification
// @Description Update a notification
// @Tags notification
// @Accept  json
// @Produce  json
// @Param id path int true "Notification ID"
// @Param notification body Request true "Notification"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /notification/{id} [put]
func UpdateNotification(w http.ResponseWriter, r *http.Request) {
	n := utils.GetNotificationCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	expireAt, err := http.ParseTime(data.ExpireAt)
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	// Make sure expireAt is in the future
	if expireAt.Before(time.Now()) {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("expire_at must be in the future")))
		return
	}

	if !models.IsValidUser(data.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	n.CID = data.CID
	n.Category = data.Category
	n.Title = data.Title
	n.Body = data.Body
	n.ExpireAt = expireAt

	if err := n.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewNotificationResponse(n))
}

// PatchNotification godoc
// @Summary Patch a notification
// @Description Patch a notification
// @Tags notification
// @Accept  json
// @Produce  json
// @Param id path int true "Notification ID"
// @Param notification body Request true "Notification"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /notification/{id} [patch]
func PatchNotification(w http.ResponseWriter, r *http.Request) {
	n := utils.GetNotificationCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.CID != 0 {
		if !models.IsValidUser(data.CID) {
			utils.Render(w, r, utils.ErrInvalidCID)
			return
		}
		n.CID = data.CID
	}
	if data.Category != "" {
		n.Category = data.Category
	}
	if data.Title != "" {
		n.Title = data.Title
	}
	if data.Body != "" {
		n.Body = data.Body
	}
	if data.ExpireAt != "" {
		expireAt, err := http.ParseTime(data.ExpireAt)
		if err != nil {
			utils.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}

		// Make sure expireAt is in the future
		if expireAt.Before(time.Now()) {
			utils.Render(w, r, utils.ErrInvalidRequest(errors.New("expire_at must be in the future")))
			return
		}

		n.ExpireAt = expireAt
	}

	if err := n.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewNotificationResponse(n))
}

// DeleteNotification godoc
// @Summary Delete a notification
// @Description Delete a notification
// @Tags notification
// @Accept  json
// @Produce  json
// @Param id path int true "Notification ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /notification/{id} [delete]
func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	n := utils.GetNotificationCtx(r)
	if err := n.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
}
