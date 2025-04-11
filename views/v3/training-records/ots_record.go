package trainingrecords

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

type otsRecordRequest struct {
	Data   datatypes.JSON `json:"data"`
	Notes  string         `json:"notes"`
	Result bool           `json:"result"`
}

func (req *otsRecordRequest) Validate() error {
	return validator.New().Struct(req)
}

func (req *otsRecordRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type otsRecordResponse struct {
	*models.RatingExamRecord
}

func newOTSRecordResponse(temp *models.RatingExamRecord) *otsRecordResponse {
	return &otsRecordResponse{temp}
}

func (res *otsRecordResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.RatingExamRecord == nil {
		return errors.New("OTS Record not found")
	}
	return nil
}

func newOTSRecordList(templates []models.RatingExamRecord) []render.Renderer {
	list := []render.Renderer{}
	for idx := range templates {
		list = append(list, newOTSRecordResponse(&templates[idx]))
	}
	return list
}

// createOTSRecord godoc
// @Summary Create a new OTS Record
// @Description Create a new OTS Record
// @Tags training
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param otsRecordRequest body otsRecordRequest true "OTS Record"
// @Success 201 {object} otsRecordResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/ots [post]
func createOTSRecord(w http.ResponseWriter, r *http.Request) {
	cid := utils.GetUserCtx(r).CID
	data := &otsRecordRequest{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	record := &models.RatingExamRecord{
		StudentCID: cid,
		Data:       data.Data,
		Notes:      data.Notes,
		Result:     data.Result,
	}

	instructorCID := utils.GetXUser(r).CID
	record.InstructorCID = instructorCID

	if err := record.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.Render(w, r, newOTSRecordResponse(record)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// listOTSRecords godoc
// @Summary List all OTS Records for user
// @Description List all OTS Records for user
// @Tags training
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Success 200 {object} []otsRecordResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/ots [get]
func listOTSRecords(w http.ResponseWriter, r *http.Request) {
	cid := utils.GetUserCtx(r).CID
	records, err := models.GetFilteredOTSRecords(map[string]interface{}{"student_cid": cid})
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	if err := render.RenderList(w, r, newOTSRecordList(records)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}
