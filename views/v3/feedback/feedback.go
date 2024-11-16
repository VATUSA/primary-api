package feedback

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	PilotCID      uint                 `json:"pilot_cid" example:"1293257" validate:"required"`
	Callsign      string               `json:"callsign" example:"DAL123" validate:"required"`
	ControllerCID uint                 `json:"controller_cid" example:"1293257" validate:"required"`
	Position      string               `json:"position" example:"DEN_I_APP" validate:"required"`
	Rating        types.FeedbackRating `json:"rating" example:"good" validate:"required,oneof=unsatisfactory poor fair good excellent"`
	Feedback      string               `json:"feedback" example:"Raaj was the best controller I've ever flown under." validate:"required"`
	Status        types.StatusType     `json:"status" example:"pending" validate:"required,oneof=pending approved denied"`
	Comment       string               `json:"comment" example:"Great work Raaj!"`
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
	ID                  uint                 `json:"id" example:"1"`
	ControllerCID       uint                 `json:"controller_cid" example:"1293257"`
	ControllerFirstName string               `json:"controller_first_name" example:"John"`
	ControllerLastName  string               `json:"controller_last_name" example:"Doe"`
	Callsign            string               `json:"callsign" example:"DAL123"`
	Position            string               `json:"position" example:"DEN_I_APP"`
	Facility            constants.FacilityID `json:"facility" example:"ZDV"`
	Rating              types.FeedbackRating `json:"rating" example:"good"`
	Feedback            string               `json:"feedback" example:"Raaj was the best controller I've ever flown under."`
	Status              types.StatusType     `json:"status" example:"pending"`
	Comment             string               `json:"comment" example:"Great work Raaj!"`
	PilotCID            uint                 `json:"pilot_cid" example:"1293257"`
	CreatedAt           time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
}

func NewFeedbackResponse(f *models.Feedback) *Response {
	resp := &Response{
		ID:            f.ID,
		ControllerCID: f.ControllerCID,
		Callsign:      f.Callsign,
		Position:      f.Position,
		Facility:      f.Facility,
		Rating:        f.Rating,
		Feedback:      f.Feedback,
		Status:        f.Status,
		Comment:       f.Comment,
		PilotCID:      f.PilotCID,
		CreatedAt:     f.CreatedAt,
	}

	controller := &models.User{CID: f.ControllerCID}
	if err := controller.Get(); err == nil {
		resp.ControllerFirstName = controller.FirstName
		resp.ControllerLastName = controller.LastName
	}

	return resp
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.ID == 0 {
		return errors.New("feedback is required")
	}
	return nil
}

func NewFeedbackListResponse(feedback []models.Feedback) []render.Renderer {
	list := []render.Renderer{}
	for idx := range feedback {
		list = append(list, NewFeedbackResponse(&feedback[idx]))
	}
	return list
}

// CreateFeedback godoc
// @Summary Create a new feedback entry
// @Description Create a new feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param facility path string true "Facility ID"
// @Param feedback body Request true "Feedback Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/feedback [post]
func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(data.ControllerCID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !models.IsValidFacility(fac.ID) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	f := &models.Feedback{
		PilotCID:      data.PilotCID,
		Callsign:      data.Callsign,
		ControllerCID: data.ControllerCID,
		Position:      data.Position,
		Facility:      fac.ID,
		Rating:        data.Rating,
		Feedback:      data.Feedback,
		Status:        data.Status,
		Comment:       data.Comment,
	}

	if f.Status == types.Pending {
		f.Comment = ""
	}

	if err := f.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewFeedbackResponse(f))
}

// ListFeedback godoc
// @Summary List feedback entries
// @Description List feedback entries
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param facility path string true "Facility ID"
// @Param cid query int false "CID"
// @Param status query string false "Status"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/feedback [get]
func ListFeedback(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("cid")

	status := types.StatusType(r.URL.Query().Get("status"))

	if cid != "" {
		cidInt, err := strconv.Atoi(cid)
		if err != nil {
			utils.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}

		f, err := models.GetFeedbackByFacilityAndCID(utils.GetFacilityCtx(r).ID, uint(cidInt), status)
		if err != nil {
			utils.Render(w, r, utils.ErrInternalServer)
			return
		}

		if err := render.RenderList(w, r, NewFeedbackListResponse(f)); err != nil {
			utils.Render(w, r, utils.ErrRender(err))
			return
		}
		return
	}

	f, err := models.GetFeedbackByFacility(utils.GetFacilityCtx(r).ID, status)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewFeedbackListResponse(f)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateFeedback godoc
// @Summary Update a feedback entry
// @Description Update a feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param id path int true "Feedback ID"
// @Param facility path string true "Facility ID"
// @Param feedback body Request true "Feedback Entry"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/feedback/{id} [put]
func UpdateFeedback(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(data.ControllerCID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	f := utils.GetFeedbackCtx(r)
	f.PilotCID = data.PilotCID
	f.Callsign = data.Callsign
	f.ControllerCID = data.ControllerCID
	f.Position = data.Position
	f.Rating = data.Rating
	f.Feedback = data.Feedback
	f.Status = data.Status
	f.Comment = data.Comment

	if err := f.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// PatchFeedback godoc
// @Summary Patch a feedback entry
// @Description Patch a feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param id path int true "Feedback ID"
// @Param facility path string true "Facility ID"
// @Param feedback body Request true "Feedback Entry"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/feedback/{id} [patch]
func PatchFeedback(w http.ResponseWriter, r *http.Request) {
	f := utils.GetFeedbackCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.PilotCID != 0 {
		f.PilotCID = data.PilotCID
	}
	if data.Callsign != "" {
		f.Callsign = data.Callsign
	}
	if data.ControllerCID != 0 {
		if !models.IsValidUser(data.ControllerCID) {
			utils.Render(w, r, utils.ErrInvalidCID)
			return
		}
		f.ControllerCID = data.ControllerCID
	}
	if data.Position != "" {
		f.Position = data.Position
	}
	if data.Rating != "" {
		f.Rating = data.Rating
	}
	if data.Feedback != "" {
		f.Feedback = data.Feedback
	}
	if data.Status != "" {
		f.Status = data.Status
	}
	if data.Comment != "" {
		f.Comment = data.Comment
	}

	if err := f.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// DeleteFeedback godoc
// @Summary Delete a feedback entry
// @Description Delete a feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param id path int true "Feedback ID"
// @Param facility path string true "Facility ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/feedback/{id} [delete]
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	f := utils.GetFeedbackCtx(r)
	if err := f.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// GetUserFeedback godoc
// @Summary Get accepted feedback entries for a user
// @Description Get feedback entries for a user
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param CID path int true "CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/feedback [get]
func GetUserFeedback(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	status := types.Accepted

	feedbacks, err := models.GetFeedbackByCID(user.CID, status)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewFeedbackListResponse(feedbacks)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}
