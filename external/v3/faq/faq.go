package faq

import (
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Question string `json:"question" validate:"required" example:"What ARTCC should I join?"`
	Answer   string `json:"answer" validate:"required" example:"You should join ZDV."`
	Category string `json:"category" validate:"required,oneof=membership training technology misc" example:"membership"`
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
	for idx := range faqs {
		list = append(list, NewFAQResponse(&faqs[idx]))
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
// @Param FacilityID path string true "Facility ID"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/faq [post]
func CreateFAQ(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if !models.IsValidFacility(fac.ID) {
		utils.Render(w, r, utils.ErrInvalidFacility)
		return
	}

	faq := &models.FAQ{
		Facility:  fac.ID,
		Question:  data.Question,
		Answer:    data.Answer,
		Category:  data.Category,
		CreatedBy: 111,
	}

	if self := utils.GetXUser(r); self != nil {
		faq.CreatedBy = self.CID
	}

	if err := faq.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewFAQResponse(faq))

}

// ListFAQ godoc
// @Summary List all FAQs
// @Description List all FAQs
// @Tags faq
// @Accept  json
// @Produce  json
// @Param FacilityID query string false "Facility ID"
// @Success 200 {object} []Response
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/faq [get]
func ListFAQ(w http.ResponseWriter, r *http.Request) {
	fac := utils.GetFacilityCtx(r)

	faqs, err := models.GetAllFAQByFacility(fac.ID)
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
// @Param FacilityID path string true "Facility ID"
// @Param faq body Request true "FAQ"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/faq/{id} [put]
func UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	faq := utils.GetFAQCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	faq.Question = data.Question
	faq.Answer = data.Answer
	faq.Category = data.Category

	if self := utils.GetXUser(r); self != nil {
		faq.UpdatedBy = self.CID
	}

	if err := faq.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFAQResponse(faq))
}

// PatchFAQ godoc
// @Summary Patch a FAQ
// @Description Patch a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Param FacilityID path string true "Facility ID"
// @Param faq body Request true "FAQ"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/faq/{id} [patch]
func PatchFAQ(w http.ResponseWriter, r *http.Request) {
	faq := utils.GetFAQCtx(r)

	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
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

	if self := utils.GetXUser(r); self != nil {
		faq.UpdatedBy = self.CID
	}

	if err := faq.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewFAQResponse(faq))
}

// DeleteFAQ godoc
// @Summary Delete a FAQ
// @Description Delete a FAQ
// @Tags faq
// @Accept  json
// @Produce  json
// @Param id path string true "FAQ ID"
// @Param FacilityID path string true "Facility ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /facility/{FacilityID}/faq/{id} [delete]
func DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	faq := utils.GetFAQCtx(r)

	if err := faq.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
