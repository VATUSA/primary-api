package facility

import (
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Name       string `json:"name" example:"Seattle ARTCC" validate:"required"`
	About      string `json:"about" example:"Seattle ARTCC contains ZSE... etc. etc. etc." validate:"required"`
	URL        string `json:"url" example:"https://zseartcc.org" validate:"required"`
	WebhookURL string `json:"webhook_url" example:"" validate:"required"`
}

func (req *Request) Validate() error {
	return nil
}

func (req *Request) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type Response struct {
	ID         constants.FacilityID `json:"id" example:"ZDV"`
	Name       string               `json:"name" example:"Denver ARTCC"`
	About      string               `json:"about" example:"Denver ARTCC contains ZDV... etc. etc. etc."`
	URL        string               `json:"url" example:"https://zdvartcc.org"`
	WebhookURL string               `json:"webhook_url" example:""`
}

func NewFacilityResponse(facility *models.Facility) *Response {
	resp := &Response{
		ID:         facility.ID,
		About:      facility.About,
		Name:       facility.Name,
		URL:        facility.URL,
		WebhookURL: facility.WebhookURL,
	}

	return resp
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewFacilityListResponse(facilities []models.Facility) []render.Renderer {
	list := []render.Renderer{}
	for idx := range facilities {
		list = append(list, NewFacilityResponse(&facilities[idx]))
	}
	return list
}

// GetFacilities godoc
// @Summary Get all facilities
// @Description Get all facilities
// @Tags facility
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility [get]
func GetFacilities(w http.ResponseWriter, r *http.Request) {
	facs, err := models.GetAllFacilities()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	// Strip sensitive data
	for idx := range facs {
		facs[idx].StripSensitive()
	}

	if err := render.RenderList(w, r, NewFacilityListResponse(facs)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// GetFacility godoc
// @Summary Get a specific facility
// @Description Get information for a specific facility
// @Tags facility
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID} [get]
func GetFacility(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	utils.Render(w, r, NewFacilityResponse(fac))
}

// UpdateFacility godoc
// @Summary Update a facility
// @Description Update a facility
// @Tags facility
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param facility body Request true "Facility"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID} [put]
func UpdateFacility(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	fac.Name = req.Name
	fac.About = req.About
	fac.URL = req.URL
	fac.WebhookURL = req.WebhookURL

	if err := fac.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFacilityResponse(fac))
}

// PatchFacility godoc
// @Summary Patch a facility
// @Description Patch a facility
// @Tags facility
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Param facility body Request true "Facility"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID} [patch]
func PatchFacility(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrBadRequest)
		return
	}

	if fac.Name != "" {
		fac.Name = req.Name
	}
	if fac.About != "" {
		fac.About = req.About
	}
	if fac.URL != "" {
		fac.URL = req.URL
	}
	if fac.WebhookURL != "" {
		fac.WebhookURL = req.WebhookURL
	}

	if err := fac.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFacilityResponse(fac))
}

// ResetApiKey godoc
// @Summary Regenerate an API key
// @Description Regenerate an API key
// @Tags facility
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/reset-api-key [post]
func ResetApiKey(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	key, err := models.GenerateApiKey()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	fac.APIKey = key
	if err := fac.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.JSON(w, r, http.StatusOK, map[string]string{"api_key": key})
}
