package roster_request

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Facility    string            `json:"requested_facility" example:"ZDV" validate:"required,len=3"`
	RequestType types.RequestType `json:"request_type" example:"visiting" validate:"required,oneof=visiting transferring"`
	Status      types.StatusType  `json:"status" example:"pending" validate:"required,oneof=pending accepted rejected"`
	Reason      string            `json:"reason" example:"I want to transfer to ZDV" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.RosterRequest
}

func NewRosterRequestResponse(r *models.RosterRequest) *Response {
	return &Response{RosterRequest: r}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.RosterRequest == nil {
		return errors.New("roster request not found")
	}
	return nil
}

func NewRosterRequestListResponse(r []models.RosterRequest) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range r {
		list = append(list, NewRosterRequestResponse(&d))
	}
	return list
}

// CreateRosterRequest godoc
// @Summary Create a new roster request
// @Description Create a new roster request
// @Tags roster-request
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param roster_request body Request true "Roster Request"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster-request [post]
func CreateRosterRequest(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := req.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := utils.GetUserCtx(r)

	if !models.IsValidUser(user.CID) {
		render.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !models.IsValidFacility(req.Facility) {
		render.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	rosterRequest := &models.RosterRequest{
		CID:         user.CID,
		Facility:    req.Facility,
		RequestType: req.RequestType,
		Status:      req.Status,
		Reason:      req.Reason,
	}

	if err := rosterRequest.Create(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRosterRequestResponse(rosterRequest))
}

// ListRosterRequest godoc
// @Summary List all roster requests
// @Description List all roster requests
// @Tags roster-request
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster-request [get]
func ListRosterRequest(w http.ResponseWriter, r *http.Request) {
	rosterRequests, err := models.GetAllRosterRequestsByCID(utils.GetUserCtx(r).CID)
	if err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewRosterRequestListResponse(rosterRequests)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateRosterRequest godoc
// @Summary Update a roster request
// @Description Update a roster request
// @Tags roster-request
// @Accept  json
// @Produce  json
// @Param id path string true "Roster Request ID"
// @Param roster_request body Request true "Roster Request"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster-request/{id} [put]
func UpdateRosterRequest(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(data.Facility) {
		render.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	req := utils.GetRosterRequestCtx(r)

	if req.Status == types.Pending && data.Status == types.Accepted {
		roster := &models.Roster{
			CID:      req.CID,
			Facility: data.Facility,
			OIs:      "",
			Home:     false,
			Visiting: false,
			Status:   "Active",
		}

		if data.RequestType == types.Visiting {
			roster.Visiting = true
		} else {
			roster.Home = true
		}

		if err := roster.Create(); err != nil {
			render.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}
	}

	req.Facility = data.Facility
	req.RequestType = data.RequestType
	req.Status = data.Status
	req.Reason = data.Reason

	if err := req.Update(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewRosterRequestResponse(req))
}

// DeleteRosterRequest godoc
// @Summary Delete a roster request
// @Description Delete a roster request
// @Tags roster-request
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param id path string true "Roster Request ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster-request/{id} [delete]
func DeleteRosterRequest(w http.ResponseWriter, r *http.Request) {
	req := utils.GetRosterRequestCtx(r)
	if err := req.Delete(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
