package facility

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Name string `json:"name" example:"Seattle ARTCC" validate:"required"`
	URL  string `json:"url" example:"https://zseartcc.org" validate:"required"`
}

func (req *Request) Validate() error {
	return nil
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.Facility
}

func NewFacilityResponse(facility *models.Facility) *Response {
	resp := &Response{Facility: facility}

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
	fac.URL = req.URL

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
	if fac.URL != "" {
		fac.URL = req.URL
	}

	if err := fac.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFacilityResponse(fac))
}
