package feedback

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	PilotCID      uint                 `json:"pilot_cid" example:"1293257" validate:"required"`
	Callsign      string               `json:"callsign" example:"DAL123" validate:"required"`
	ControllerCID uint                 `json:"controller_cid" example:"1293257" validate:"required"`
	Position      string               `json:"position" example:"DEN_I_APP" validate:"required"`
	Facility      constants.FacilityID `json:"facility" example:"ZDV" validate:"required,len=3"`
	Rating        types.FeedbackRating `json:"rating" example:"good" validate:"required,oneof=unsatisfactory poor fair good excellent"`
	Notes         string               `json:"notes" example:"Raaj was the best controller I've ever flown under." validate:"required"`
	Status        types.StatusType     `json:"status" example:"pending" validate:"required,oneof=pending approved denied"`
	Comment       string               `json:"comment" example:"Great work Raaj!"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.Feedback
}

func NewFeedbackResponse(f *models.Feedback) *Response {
	return &Response{Feedback: f}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Feedback == nil {
		return errors.New("feedback not found")
	}
	return nil
}

func NewFeedbackListResponse(f []models.Feedback) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range f {
		list = append(list, NewFeedbackResponse(&d))
	}
	return list
}

// CreateFeedback godoc
// @Summary Create a new feedback entry
// @Description Create a new feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param feedback body Request true "Feedback Entry"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback [post]
func CreateFeedback(w http.ResponseWriter, r *http.Request) {
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

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	f := &models.Feedback{
		PilotCID:      data.PilotCID,
		Callsign:      data.Callsign,
		ControllerCID: data.ControllerCID,
		Position:      data.Position,
		Facility:      data.Facility,
		Rating:        data.Rating,
		Notes:         data.Notes,
		Status:        data.Status,
		Comment:       data.Comment,
	}
	if err := f.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewFeedbackResponse(f))
}

// GetFeedback godoc
// @Summary Get a feedback entry
// @Description Get a feedback entry
// @Tags feedback
// @Accept  json
// @Produce  json
// @Param id path int true "Feedback ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback/{id} [get]
func GetFeedback(w http.ResponseWriter, r *http.Request) {
	f := utils.GetFeedbackCtx(r)
	utils.Render(w, r, NewFeedbackResponse(f))
}

// ListFeedback godoc
// @Summary List feedback entries
// @Description List feedback entries
// @Tags feedback
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback [get]
func ListFeedback(w http.ResponseWriter, r *http.Request) {
	f, err := models.GetAllFeedback()
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
// @Param feedback body Request true "Feedback Entry"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback/{id} [put]
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

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	f := utils.GetFeedbackCtx(r)
	f.PilotCID = data.PilotCID
	f.Callsign = data.Callsign
	f.ControllerCID = data.ControllerCID
	f.Position = data.Position
	f.Facility = data.Facility
	f.Rating = data.Rating
	f.Notes = data.Notes
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
// @Param feedback body Request true "Feedback Entry"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback/{id} [patch]
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
	if data.Facility != "" {
		if !models.IsValidFacility(data.Facility) {
			utils.Render(w, r, utils.ErrInvalidFacility)
			return
		}
		f.Facility = data.Facility
	}
	if data.Rating != "" {
		f.Rating = data.Rating
	}
	if data.Notes != "" {
		f.Notes = data.Notes
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
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /feedback/{id} [delete]
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	f := utils.GetFeedbackCtx(r)
	if err := f.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
