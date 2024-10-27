package action_log

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
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
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
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
	for idx := range ale {
		list = append(list, NewActionLogEntryResponse(&ale[idx]))
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

	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	ale := &models.ActionLogEntry{
		CID:       user.CID,
		Entry:     data.Entry,
		CreatedBy: "System",
	}

	if err := ale.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewActionLogEntryResponse(ale))
}

// ListActionLog godoc
// @Summary List all action log entries
// @Description List all action log entries
// @Tags action-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/action-log [get]
func ListActionLog(w http.ResponseWriter, r *http.Request) {
	ale, err := models.GetAllActionLogEntriesByCID(utils.GetUserCtx(r).CID)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewActionLogEntryListResponse(ale)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateActionLog godoc
// @Summary Update an action log entry
// @Description Update an action log entry
// @Tags action-log
// @Accept  json
// @Produce  json
// @Param id path int true "Action Log Entry ID"
// @Param action_log body Request true "Action Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/action-log/{id} [put]
func UpdateActionLog(w http.ResponseWriter, r *http.Request) {
	ale := utils.GetActionLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	ale.Entry = data.Entry
	ale.UpdatedBy = "System"

	if err := ale.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewActionLogEntryResponse(ale))
}

// PatchActionLog godoc
// @Summary Patch an action log entry
// @Description Patch an action log entry
// @Tags action-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Action Log Entry ID"
// @Param action_log body Request true "Action Log Entry"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/action-log/{id} [patch]
func PatchActionLog(w http.ResponseWriter, r *http.Request) {
	ale := utils.GetActionLogCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Entry != "" {
		ale.Entry = data.Entry
	}
	ale.UpdatedBy = "System"

	if err := ale.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewActionLogEntryResponse(ale))
}

// DeleteActionLog godoc
// @Summary Delete an action log entry
// @Description Delete an action log entry
// @Tags action-log
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Action Log Entry ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/action-log/{id} [delete]
func DeleteActionLog(w http.ResponseWriter, r *http.Request) {
	ale := utils.GetActionLogCtx(r)

	if err := ale.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
