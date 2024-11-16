package roster

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Request struct {
	CID      uint   `json:"cid" example:"123456" validate:"required"`
	OIs      string `json:"operating_initials" example:"RP" validate:"required"`
	Home     bool   `json:"home" example:"true"`
	Visiting bool   `json:"visiting" example:"false"`
	Status   string `json:"status" example:"Active" validate:"required,oneof=active loa"` // Active, LOA
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
	*models.Roster
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewRosterResponse(r *models.Roster) *Response {
	// TODO - Speed this up by putting the user in the roster struct and sending a custom struct response
	u := &models.User{CID: r.CID}
	if err := u.Get(); err != nil {
		return &Response{Roster: r}
	}

	return &Response{
		Roster:    r,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Roster == nil {
		return errors.New("roster not found")
	}

	return nil
}

func NewRosterListResponse(rosters []models.Roster) []render.Renderer {
	list := []render.Renderer{}
	for idx := range rosters {
		list = append(list, NewRosterResponse(&rosters[idx]))
	}

	return list
}

// CreateRoster godoc
// @Summary Add user to roster
// @Description Adds a user to a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param roster body Request true "Roster"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/roster [post]
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

	if !models.IsValidUser(data.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	fac := utils.GetFacilityCtx(r)

	if !data.Home && !data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be false")))
		return
	}

	if data.Home && data.Visiting {
		utils.Render(w, r, utils.ErrInvalidRequest(errors.New("home and visiting cannot both be true")))
		return
	}

	roster := &models.Roster{
		CID:      data.CID,
		Facility: fac.ID,
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

// GetRosterByFacility godoc
// @Summary Get rosters by facility
// @Description Get rosters by facility
// @Tags roster
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param type query string false "Type" Enums(home,visiting)
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/roster [get]
func GetRosterByFacility(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	rosterType := r.URL.Query().Get("type")

	if rosterType == "" {
		rosters, err := models.GetRostersByFacility(fac.ID)
		if err != nil {
			utils.Render(w, r, utils.ErrInternalServer)
			return
		}

		if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
			utils.Render(w, r, utils.ErrRender(err))
			return
		}
		return
	}

	rosters, err := models.GetRostersByFacilityAndType(fac.ID, strings.ToLower(rosterType))
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewRosterListResponse(rosters)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// DeleteRoster godoc
// @Summary Delete a roster
// @Description Delete a roster
// @Tags roster
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param id path int true "Roster ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/roster/{id} [delete]
func DeleteRoster(w http.ResponseWriter, r *http.Request) {
	roster := utils.GetRosterCtx(r)

	if err := roster.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}

// GetUserRosters godoc
// @Summary Get rosters by user
// @Description Get rosters by user
// @Tags roster
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/roster [get]
func GetUserRosters(w http.ResponseWriter, r *http.Request) {
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
