package roster

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Request struct {
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
// @Param cid path int true "User CID"
// @Param roster body Request true "Roster"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster [post]
func CreateRoster(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
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

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	if !data.Home && !data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be false")))
		return
	}

	if data.Home && data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be true")))
		return
	}

	roster := &models.Roster{
		CID:      user.CID,
		Facility: data.Facility,
		OIs:      data.OIs,
		Home:     data.Home,
		Visiting: data.Visiting,
		Status:   data.Status,
	}

	if err := roster.Create(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewRosterResponse(roster))
}

// GetRoster godoc
// @Summary Get a roster
// @Description Get a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Roster ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster/{id} [get]
func GetRoster(w http.ResponseWriter, r *http.Request) {
	roster := utils.GetRosterCtx(r)

	utils.Render(w, r, NewRosterResponse(roster))
}

// ListRoster godoc
// @Summary List rosters for a user
// @Description List rosters for a user
// @Tags roster
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster [get]
func ListRoster(w http.ResponseWriter, r *http.Request) {
	user := utils.GetUserCtx(r)

	rosters, err := models.GetRostersByCID(user.CID)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateRoster godoc
// @Summary Update a roster
// @Description Update a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param cid path int true "User CID"
// @Param id path int true "Roster ID"
// @Param roster body Request true "Roster"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster/{id} [put]
func UpdateRoster(w http.ResponseWriter, r *http.Request) {
	roster := utils.GetRosterCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
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

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	if !data.Home && !data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be false")))
		return
	}

	if data.Home && data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be true")))
		return
	}

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
// @Param cid path int true "User CID"
// @Param id path int true "Roster ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster/{id} [delete]
func DeleteRoster(w http.ResponseWriter, r *http.Request) {
	roster := utils.GetRosterCtx(r)

	if err := roster.Delete(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// GetRosterByFacility godoc
// @Summary Get rosters by facility
// @Description Get rosters by facility
// @Tags roster
// @Accept  json
// @Produce  json
// @Param facility path string true "Facility"
// @Param type query string false "Type" Enums(home,visiting)
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{facility}/roster [get]
func GetRosterByFacility(w http.ResponseWriter, r *http.Request) {
	fac, err := utils.GetFacilityCtx(r)
	if err != nil {
		render.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	rosterType := r.URL.Query().Get("type")

	if rosterType == "" {
		rosters, err := models.GetRostersByFacility(fac.ID)
		if err != nil {
			render.Render(w, r, utils.ErrInternalServer)
			return
		}

		if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
			render.Render(w, r, utils.ErrRender(err))
			return
		}
		return
	}

	rosters, err := models.GetRostersByFacilityAndType(fac.ID, strings.ToLower(rosterType))
	if err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}
