package user_flag

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	NoStaffRole              bool `json:"no_staff_role" example:"false"`
	NoStaffLogEntryID        uint `json:"no_staff_log_entry_id" example:"1"`
	NoVisiting               bool `json:"no_visiting" example:"false"`
	NoVisitingLogEntryID     uint `json:"no_visiting_log_entry_id" example:"1"`
	NoTransferring           bool `json:"no_transferring" example:"false"`
	NoTransferringLogEntryID uint `json:"no_transferring_log_entry_id" example:"`
	NoTraining               bool `json:"no_training" example:"false"`
	NoTrainingLogEntryID     uint `json:"no_training_log_entry_id" example:"1"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.UserFlag
}

func NewUserFlagResponse(r *models.UserFlag) *Response {
	return &Response{UserFlag: r}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.UserFlag == nil {
		return errors.New("user flag not found")
	}

	return nil
}

func NewUserFlagListResponse(r []models.UserFlag) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range r {
		list = append(list, NewUserFlagResponse(&d))
	}

	return list
}

// GetUserFlag godoc
// @Summary Get a user flag
// @Description Get a user flag
// @Tags user-flag
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/user-flag [get]
func GetUserFlag(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewUserFlagResponse(utils.GetUserFlagCtx(r)))
}

// UpdateUserFlag godoc
// @Summary Update a user flag
// @Description Update a user flag
// @Tags user-flag
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param user_flag body Request true "User Flag"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/user-flag [put]
func UpdateUserFlag(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	userFlag := utils.GetUserFlagCtx(r)
	userFlag.NoStaffRole = req.NoStaffRole
	userFlag.NoStaffLogEntryID = req.NoStaffLogEntryID
	userFlag.NoVisiting = req.NoVisiting
	userFlag.NoVisitingLogEntryID = req.NoVisitingLogEntryID
	userFlag.NoTransferring = req.NoTransferring
	userFlag.NoTransferringLogEntryID = req.NoTransferringLogEntryID
	userFlag.NoTraining = req.NoTraining
	userFlag.NoTrainingLogEntryID = req.NoTrainingLogEntryID

	if err := userFlag.Update(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserFlagResponse(userFlag))
}

// PatchUserFlag godoc
// @Summary Patch a user flag
// @Description Patch a user flag
// @Tags user-flag
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param user_flag body Request true "User Flag"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/user-flag [patch]
func PatchUserFlag(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	userFlag := utils.GetUserFlagCtx(r)
	if req.NoStaffRole {
		userFlag.NoStaffRole = req.NoStaffRole
	}
	if req.NoStaffLogEntryID != 0 {
		userFlag.NoStaffLogEntryID = req.NoStaffLogEntryID
	}
	if req.NoVisiting {
		userFlag.NoVisiting = req.NoVisiting
	}
	if req.NoVisitingLogEntryID != 0 {
		userFlag.NoVisitingLogEntryID = req.NoVisitingLogEntryID
	}
	if req.NoTransferring {
		userFlag.NoTransferring = req.NoTransferring
	}
	if req.NoTransferringLogEntryID != 0 {
		userFlag.NoTransferringLogEntryID = req.NoTransferringLogEntryID
	}
	if req.NoTraining {
		userFlag.NoTraining = req.NoTraining
	}
	if req.NoTrainingLogEntryID != 0 {
		userFlag.NoTrainingLogEntryID = req.NoTrainingLogEntryID
	}

	if err := userFlag.Update(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserFlagResponse(userFlag))
}

// DeleteUserFlag godoc
// @Summary Delete a user flag
// @Description Delete a user flag
// @Tags user-flag
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/user-flag [delete]
func DeleteUserFlag(w http.ResponseWriter, r *http.Request) {
	userFlag := utils.GetUserFlagCtx(r)
	if err := userFlag.Delete(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
