package training

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"net/http"
)

type otsTemplateRequest struct {
	Name     string          `json:"name"`
	Template json.RawMessage `json:"template"`
}

func (req *otsTemplateRequest) Validate() error {
	return validator.New().Struct(req)
}

func (req *otsTemplateRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type otsTemplateResponse struct {
	*models.OTSTemplate
}

func newOTSTemplateResponse(temp *models.OTSTemplate) *otsTemplateResponse {
	return &otsTemplateResponse{temp}
}

func (res *otsTemplateResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.OTSTemplate == nil {
		return errors.New("OTS Template not found")
	}
	return nil
}

func newTemplateList(templates []models.OTSTemplate) []render.Renderer {
	list := []render.Renderer{}
	for idx := range templates {
		list = append(list, newOTSTemplateResponse(&templates[idx]))
	}
	return list
}

// createTemplate godoc
// @Summary Create an OTS template
// @Description Create an OTS template
// @Tags training
// @Accept  json
// @Produce  json
// @Param template body otsTemplateRequest true "Template"
// @Success 201 {object} otsTemplateResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /training/ots/templates [post]
func createTemplate(w http.ResponseWriter, r *http.Request) {
	data := &otsTemplateRequest{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	template := &models.OTSTemplate{
		Name: data.Name,
	}

	if err := json.Unmarshal(data.Template, &template.Template); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	template.Template = datatypes.JSON(data.Template)

	user := utils.GetXUser(r)
	template.LastUpdatedBy = user.CID

	if err := template.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, newOTSTemplateResponse(template))
}

// listTemplates godoc
// @Summary List all OTS templates
// @Description List all OTS templates
// @Tags training
// @Accept  json
// @Produce  json
// @Success 200 {object} []otsTemplateResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /training/ots/templates [get]
func listTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := models.GetAllOTSTemplates()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, newTemplateList(templates)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// patchTemplate godoc
// @Summary Patch an OTS template
// @Description Patch an OTS template
// @Tags training
// @Accept  json
// @Produce  json
// @Param id path int true "Template ID"
// @Param template body otsTemplateRequest true "Template"
// @Success 200 {object} otsTemplateResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /training/ots/templates/{id} [patch]
func patchTemplate(w http.ResponseWriter, r *http.Request) {
	template := utils.GetOTSTemplateCtx(r)
	data := &otsTemplateRequest{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.Name != "" {
		template.Name = data.Name
	}
	if len(data.Template) != 0 {
		// Try to unmarshal the template
		if err := json.Unmarshal(data.Template, &template.Template); err != nil {
			utils.Render(w, r, utils.ErrInvalidRequest(err))
			return
		}

		template.Template = datatypes.JSON(data.Template)
	}

	user := utils.GetXUser(r)
	template.LastUpdatedBy = user.CID

	if err := template.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, newOTSTemplateResponse(template))
}

// deleteTemplate godoc
// @Summary Delete an OTS template
// @Description Delete an OTS template
// @Tags training
// @Accept  json
// @Produce  json
// @Param id path int true "Template ID"
// @Success 204
// @Failure 500 {object} utils.ErrResponse
// @Router /training/ots/templates/{id} [delete]
func deleteTemplate(w http.ResponseWriter, r *http.Request) {
	template := utils.GetOTSTemplateCtx(r)
	if err := template.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
