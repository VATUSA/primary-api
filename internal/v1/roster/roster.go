package roster

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	CID      uint   `json:"cid" example:"1293257" validate:"required"`
	Facility string `json:"facility" example:"ZDV" validate:"required,len=3"`
	OIs      string `json:"operating_initials" example:"RP" validate:"required"`
	Home     bool   `json:"home" example:"true"`
	Visiting bool   `json:"visiting" example:"false"`
	Status   string `json:"status" example:"Active" validate:"required,oneof=active loa"` // Active, LOA
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.Roster
}

func NewRosterResponse(r *models.Roster) *Response {
	return &Response{Roster: r}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Roster == nil {
		return errors.New("roster not found")
	}

	return nil
}

func NewRosterListResponse(r []models.Roster) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range r {
		list = append(list, NewRosterResponse(&d))
	}

	return list
}

// CreateRoster godoc
// @Summary Create a new roster
// @Description Create a new roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param roster body Request true "Roster"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /roster [post]
func CreateRoster(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(data.CID) {
		render.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !models.IsValidFacility(data.Facility) {
		render.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	if !data.Home && !data.Visiting {
		render.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be false")))
		return
	}

	if data.Home && data.Visiting {
		render.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be true")))
		return
	}

	roster := &models.Roster{
		CID:      data.CID,
		Facility: data.Facility,
		OIs:      data.OIs,
		Home:     data.Home,
		Visiting: data.Visiting,
		Status:   data.Status,
	}

	if err := roster.Create(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewRosterResponse(roster))
}

// GetRoster godoc
// @Summary Get a roster
// @Description Get a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param id path int true "Roster ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /roster/{id} [get]
func GetRoster(w http.ResponseWriter, r *http.Request) {
	roster := GetRosterCtx(r)

	render.Render(w, r, NewRosterResponse(roster))
}

// ListRoster godoc
// @Summary List rosters
// @Description List rosters
// @Tags roster
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /roster [get]
func ListRoster(w http.ResponseWriter, r *http.Request) {
	rosters, err := models.GetAllRosters()
	if err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateRoster godoc
// @Summary Update a roster
// @Description Update a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param id path int true "Roster ID"
// @Param roster body Request true "Roster"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /roster/{id} [put]
func UpdateRoster(w http.ResponseWriter, r *http.Request) {
	roster := GetRosterCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidUser(data.CID) {
		render.Render(w, r, utils.ErrInvalidCID)
		return
	}

	if !models.IsValidFacility(data.Facility) {
		render.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	if !data.Home && !data.Visiting {
		render.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be false")))
		return
	}

	if data.Home && data.Visiting {
		render.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be true")))
		return
	}

	roster.CID = data.CID
	roster.Facility = data.Facility
	roster.OIs = data.OIs
	roster.Home = data.Home
	roster.Visiting = data.Visiting
	roster.Status = data.Status

	if err := roster.Update(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// DeleteRoster godoc
// @Summary Delete a roster
// @Description Delete a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param id path int true "Roster ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /roster/{id} [delete]
func DeleteRoster(w http.ResponseWriter, r *http.Request) {
	roster := GetRosterCtx(r)

	if err := roster.Delete(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
