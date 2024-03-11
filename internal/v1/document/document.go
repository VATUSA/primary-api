package document

import (
	"errors"
	"fmt"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"github.com/VATUSA/primary-api/pkg/storage"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"path"
	"strings"
)

type Request struct {
	Facility    string `json:"facility" example:"ZDV" validate:"required,len=3"`
	Name        string `json:"name" example:"DP001" validate:"required"`
	Description string `json:"description" example:"General Division Policy" validate:"required"`
	Category    string `json:"category" example:"general" validate:"required,oneof=general training information_technology sops loas misc"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.Document
}

func NewDocumentResponse(d *models.Document) *Response {
	return &Response{Document: d}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.Document == nil {
		return errors.New("missing required document")
	}
	return nil
}

func NewDocumentListResponse(d []models.Document) []render.Renderer {
	list := []render.Renderer{}
	for _, doc := range d {
		list = append(list, NewDocumentResponse(&doc))
	}
	return list
}

// CreateDocument godoc
// @Summary Create a new document
// @Description Create a new document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param document body Request true "Document"
// @Param file formData file false "Document file"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents [post]
func CreateDocument(w http.ResponseWriter, r *http.Request, endpoint string) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	currentDocs, err := models.GetAllDocumentsByFacilityAndCategory(data.Facility, types.DocumentCategory(data.Category))
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	for _, doc := range currentDocs {
		if doc.Name == data.Name {
			utils.Render(w, r, utils.ErrInvalidRequest(fmt.Errorf("document with name %s already exists", data.Name)))
			return
		}
	}

	// Read the file from the request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	extension := path.Ext(fileHeader.Filename)
	directory := path.Join(data.Facility, data.Category)
	filename := strings.ReplaceAll(data.Name, " ", "-") + extension

	// Put the file in the S3 bucket
	if err := storage.PublicBucket.Upload(directory, filename, file); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	document := &models.Document{
		Facility:    data.Facility,
		Name:        data.Name,
		Description: data.Description,
		Category:    types.DocumentCategory(data.Category),
		URL:         path.Join(endpoint, directory, filename),
	}

	if err := document.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)

		err := storage.PublicBucket.Delete(directory, filename)
		if err != nil {
			fmt.Println("Error deleting file from S3:", err)
		}
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewDocumentResponse(document))
}

// GetDocument godoc
// @Summary Get a document
// @Description Get a document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/{id} [get]
func GetDocument(w http.ResponseWriter, r *http.Request) {
	doc := utils.GetDocumentCtx(r)
	utils.Render(w, r, NewDocumentResponse(doc))
}

// ListDocuments godoc
// @Summary List all documents
// @Description List all documents
// @Tags documents
// @Accept  json
// @Produce  json
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents [get]
func ListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := models.GetAllDocuments()
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewDocumentListResponse(docs)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// ListDocumentsByFac godoc
// @Summary List all documents by facility
// @Description List all documents by facility
// @Tags documents
// @Accept  json
// @Produce  json
// @Param Facility path string true "Facility ID"
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/facility/{Facility} [get]
func ListDocumentsByFac(w http.ResponseWriter, r *http.Request) {
	facId := chi.URLParam(r, "Facility")
	if facId == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	docs, err := models.GetAllDocumentsByFacility(facId)
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewDocumentListResponse(docs)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// ListDocumentsByFacByCat godoc
// @Summary List all documents by facility and category
// @Description List all documents by facility and category
// @Tags documents
// @Accept  json
// @Produce  json
// @Param Facility path string true "Facility ID"
// @Param Category path string true "Category"
// @Success 200 {object} []Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/facility/{Facility}/category/{Category} [get]
func ListDocumentsByFacByCat(w http.ResponseWriter, r *http.Request) {
	facId := chi.URLParam(r, "Facility")
	cat := chi.URLParam(r, "Category")
	if facId == "" || cat == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	docs, err := models.GetAllDocumentsByFacilityAndCategory(facId, types.DocumentCategory(cat))
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, NewDocumentListResponse(docs)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateDocument godoc
// @Summary Update a document
// @Description Update a document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Param document body Request true "Document"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/{id} [put]
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	doc := utils.GetDocumentCtx(r)

	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	// TODO - update the file in the S3 bucket

	doc.Facility = data.Facility
	doc.Name = data.Name
	doc.Description = data.Description
	doc.Category = types.DocumentCategory(data.Category)

	if err := doc.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDocumentResponse(doc))
}

// PatchDocument godoc
// @Summary Patch a document
// @Description Patch a document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Param document body Request true "Document"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/{id} [patch]
func PatchDocument(w http.ResponseWriter, r *http.Request) {
	doc := utils.GetDocumentCtx(r)

	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	// TODO - update the file in the S3 bucket

	if data.Facility != "" {
		doc.Facility = data.Facility
	}
	if data.Name != "" {
		doc.Name = data.Name
	}
	if data.Description != "" {
		doc.Description = data.Description
	}
	if data.Category != "" {
		doc.Category = types.DocumentCategory(data.Category)
	}

	if err := doc.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	utils.Render(w, r, NewDocumentResponse(doc))
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/{id} [delete]
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	doc := utils.GetDocumentCtx(r)

	// Delete the file from the S3 bucket
	directory := path.Join(doc.Facility, string(doc.Category))
	filename := strings.ReplaceAll(doc.Name, " ", "-") + path.Ext(doc.URL)
	if err := storage.PublicBucket.Delete(directory, filename); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := doc.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)

}

// UploadDocument godoc
// @Summary Upload a document
// @Description Upload a document
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Param file formData file true "Document file"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /documents/{id}/upload [post]
func UploadDocument(w http.ResponseWriter, r *http.Request, endpoint string) {
	data := utils.GetDocumentCtx(r)

	// Delete the file from the S3 bucket
	directory := path.Join(data.Facility, string(data.Category))
	filename := strings.ReplaceAll(data.Name, " ", "-") + path.Ext(data.URL)
	if err := storage.PublicBucket.Delete(directory, filename); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return

	}

	// Read the file from the request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	extension := path.Ext(fileHeader.Filename)
	directory = path.Join(data.Facility, string(data.Category))
	filename = strings.ReplaceAll(data.Name, " ", "-") + extension

	// Put the file in the S3 bucket
	if err := storage.PublicBucket.Upload(directory, filename, file); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	// Update the URL in the database
	data.URL = path.Join(endpoint, directory, filename)
	if err := data.Update(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
