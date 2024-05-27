package news

import (
	"errors"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Title       string `json:"title" example:"DP001 Revision 3 Released" validate:"required"`
	Description string `json:"description" example:"DP001 has been revised to include new information regarding the new VATSIM Code of Conduct" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.News
}

func NewNewsResponse(news *models.News) *Response {
	return &Response{News: news}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.News == nil {
		return errors.New("missing required news")
	}
	return nil
}

func NewNewsListResponse(news []models.News) []render.Renderer {
	list := []render.Renderer{}
	for idx := range news {
		list = append(list, NewNewsResponse(&news[idx]))
	}
	return list
}

// CreateNews godoc
// @Summary Create a new news entry
// @Description Create a new news entry
// @Tags news
// @Accept  json
// @Produce  json
// @Param news body Request true "News Entry"
// @Param FacilityID path string true "Facility ID"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/news [post]
func CreateNews(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	fac := utils.GetFacilityCtx(r)
	requestingUser := utils.GetUserCtx(r)

	news := &models.News{
		Facility:    fac.ID,
		Title:       data.Title,
		Description: data.Description,
		CreatedBy:   fmt.Sprintf("%d", requestingUser.CID),
	}

	if err := news.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewNewsResponse(news))
}

// ListNews godoc
// @Summary List all news entries
// @Description List all news entries
// @Tags news
// @Accept  json
// @Produce  json
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/news [get]
func ListNews(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	news, err := models.GetNewsByFacility(fac.ID)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewNewsListResponse(news)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateNews godoc
// @Summary Update a news entry
// @Description Update a news entry
// @Tags news
// @Accept  json
// @Produce  json
// @Param id path string true "News ID"
// @Param news body Request true "News Entry"
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/news/{id} [put]
func UpdateNews(w http.ResponseWriter, r *http.Request) {
	news := utils.GetNewsCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := req.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	requestingUser := utils.GetUserCtx(r)

	news.Title = req.Title
	news.Description = req.Description
	news.UpdatedBy = fmt.Sprintf("%d", requestingUser.CID)

	if err := news.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewNewsResponse(news))
}

// PatchNews godoc
// @Summary Patch a news entry
// @Description Patch a news entry
// @Tags news
// @Accept  json
// @Produce  json
// @Param id path string true "News ID"
// @Param news body Request true "News Entry"
// @Param FacilityID path string true "Facility ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/news/{id} [patch]
func PatchNews(w http.ResponseWriter, r *http.Request) {
	news := utils.GetNewsCtx(r)

	req := &Request{}
	if err := render.Bind(r, req); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if req.Title != "" {
		news.Title = req.Title
	}
	if req.Description != "" {
		news.Description = req.Description
	}

	requestingUser := utils.GetUserCtx(r)
	news.UpdatedBy = fmt.Sprintf("%d", requestingUser.CID)

	if err := news.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewNewsResponse(news))
}

// DeleteNews godoc
// @Summary Delete a news entry
// @Description Delete a news entry
// @Tags news
// @Accept  json
// @Produce  json
// @Param id path string true "News ID"
// @Param FacilityID path string true "Facility ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/news/{id} [delete]
func DeleteNews(w http.ResponseWriter, r *http.Request) {
	news := utils.GetNewsCtx(r)
	if err := news.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
