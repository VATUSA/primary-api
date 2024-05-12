package action_log

import (
	"fmt"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Entry string `json:"entry" example:"Changed Preferred OIs to RP" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.ActionLogEntry
}

func NewActionLogEntryResponse(ale *models.ActionLogEntry) *Response {
	return &Response{ActionLogEntry: ale}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.ActionLogEntry == nil {
		return nil
	}
	return nil
}

func NewActionLogEntryListResponse(ale []models.ActionLogEntry) []render.Renderer {
	list := []render.Renderer{}
	for _, a := range ale {
		list = append(list, NewActionLogEntryResponse(&a))
	}
	return list
}

// CreateActionLogEntry godoc
// @Summary Create a new action log entry
// @Description Create a new action log entry
// @Tags action-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param action_log body Request true "Action Log Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/action-log [post]
func CreateActionLogEntry(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := utils.GetUserCtx(r)

	self := middleware.GetSelfUser(r)

	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	ale := &models.ActionLogEntry{
		CID:       user.CID,
		Entry:     data.Entry,
		CreatedBy: fmt.Sprintf("%d", self.CID),
	}

	if err := ale.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewActionLogEntryResponse(ale))
}