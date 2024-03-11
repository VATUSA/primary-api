package faq

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Facility string `json:"facility" validate:"required,len=3"`
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	Category string `json:"category" validate:"required,oneof=membership training technology misc"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.FAQ
}

func NewFAQResponse(faq *models.FAQ) *Response {
	return &Response{FAQ: faq}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.FAQ == nil {
		return nil
	}
	return nil
}

func NewFAQListResponse(faqs []models.FAQ) []render.Renderer {
	list := []render.Renderer{}
	for _, f := range faqs {
		list = append(list, NewFAQResponse(&f))
	}
	return list
}

// CreateFAQ godoc
// @Summary Create a new FAQ
// @Description Create a new FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param faq body Request true "FAQ"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /faq [post]
func CreateFAQ(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	faq := &models.FAQ{
		Facility:  data.Facility,
		Question:  data.Question,
		Answer:    data.Answer,
		Category:  data.Category,
		CreatedBy: 1,
	}

	if err := faq.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewFAQResponse(faq))

}

// GetFAQ godoc
// @Summary Get a FAQ
// @Description Get a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /faq/{id} [get]
func GetFAQ(w http.ResponseWriter, r *http.Request) {
	faq := GetFAQCtx(r)

	utils.Render(w, r, NewFAQResponse(faq))
}

// ListFAQ godoc
// @Summary List all FAQs
// @Description List all FAQs
// @Tags faq
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 500 {object} utils.ErrResponse
// @Router /faq [get]
func ListFAQ(w http.ResponseWriter, r *http.Request) {
	faqs, err := models.GetAllFAQ()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewFAQListResponse(faqs)); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}
}

// UpdateFAQ godoc
// @Summary Update a FAQ
// @Description Update a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Param faq body Request true "FAQ"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /faq/{id} [put]
func UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	faq := GetFAQCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(data.Facility) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	faq.Facility = data.Facility
	faq.Question = data.Question
	faq.Answer = data.Answer
	faq.Category = data.Category

	if err := faq.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Render(w, r, NewFAQResponse(faq))
}

// PatchFAQ godoc
// @Summary Patch a FAQ
// @Description Patch a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Param faq body Request true "FAQ"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /faq/{id} [patch]
func PatchFAQ(w http.ResponseWriter, r *http.Request) {
	faq := GetFAQCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Facility != "" {
		if !models.IsValidFacility(data.Facility) {
			render.Render(w, r, utils.ErrInvalidFacility)
			return
		}
		faq.Facility = data.Facility
	}
	if data.Question != "" {
		faq.Question = data.Question
	}
	if data.Answer != "" {
		faq.Answer = data.Answer
	}
	if data.Category != "" {
		faq.Category = data.Category
	}

	if err := faq.Update(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Render(w, r, NewFAQResponse(faq))
}

// DeleteFAQ godoc
// @Summary Delete a FAQ
// @Description Delete a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /faq/{id} [delete]
func DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	faq := GetFAQCtx(r)

	if err := faq.Delete(); err != nil {
		render.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
