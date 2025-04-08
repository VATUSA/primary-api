package trainingrecords

import (
	"encoding/json"
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type noteRequest struct {
	Facility    constants.FacilityID `json:"facility"`
	Position    string               `json:"position"`
	Duration    time.Duration        `json:"duration"`
	Score       uint                 `json:"score"`
	Notes       string               `json:"notes"`
	OTSRecordID uint                 `json:"ots_record_id"`
	SessionDate time.Time            `json:"session_date"`
}

func (req *noteRequest) Validate() error {
	return validator.New().Struct(req)
}

func (req *noteRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type noteResponse struct {
	*models.TrainingNotes
}

func newNoteResponse(note *models.TrainingNotes) *noteResponse {
	return &noteResponse{note}
}

func (res *noteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if res.TrainingNotes == nil {
		return errors.New("training note not found")
	}
	return nil
}

func newNoteListResponse(notes []models.TrainingNotes) []render.Renderer {
	list := []render.Renderer{}
	for idx := range notes {
		list = append(list, newNoteResponse(&notes[idx]))
	}
	return list
}

// createNote godoc
// @Summary Create a new training note
// @Description Create a new training note
// @Tags training
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Param note body noteRequest true "Training Note"
// @Success 201 {object} noteResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/notes [post]
func createNote(w http.ResponseWriter, r *http.Request) {
	data := &noteRequest{}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	note := &models.TrainingNotes{
		Facility:    data.Facility,
		Position:    data.Position,
		Duration:    data.Duration,
		Score:       data.Score,
		Notes:       data.Notes,
		OTSRecordID: data.OTSRecordID,
		SessionDate: data.SessionDate,
	}

	note.StudentCID = utils.GetUserCtx(r).CID
	note.InstructorCID = utils.GetXUser(r).CID

	if err := note.Create(); err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, newNoteResponse(note))
}

// listNotes godoc
// @Summary List training notes for a user
// @Description List all training notes for a specific user
// @Tags training
// @Accept  json
// @Produce  json
// @Param cid path string true "CID"
// @Success 200 {object} noteResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/notes [get]
func listNotes(w http.ResponseWriter, r *http.Request) {
	cid := utils.GetUserCtx(r).CID

	// Fetch the training notes for the user
	notes, err := models.GetFilteredTrainingNotes(map[string]interface{}{"student_cid": cid})
	if err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	if err := render.RenderList(w, r, newNoteListResponse(notes)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// getNote godoc
// @Summary Get a specific training note
// @Description Retrieve a specific training note by its ID
// @Tags training
// @Accept  json
// @Produce  json
// @Param note_id path uint true "Note ID"
// @Param cid path string true "CID"
// @Success 200 {object} noteResponse
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/notes/{note_id} [get]
func getNote(w http.ResponseWriter, r *http.Request) {
	note := utils.GetTrainingNoteCtx(r)

	utils.Render(w, r, newNoteResponse(note))
}

// deleteNote godoc
// @Summary Delete a specific training note
// @Description Delete a specific training note by its ID
// @Tags training
// @Accept json
// @Produce json
// @Param note_id path uint true "Note ID"
// @Param cid path string true "CID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 403 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/training/notes/{note_id} [delete]
func deleteNote(w http.ResponseWriter, r *http.Request) {
	note := utils.GetTrainingNoteCtx(r)

	if err := note.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServerWithErr(err))
		return
	}

	render.Status(r, http.StatusNoContent) // HTTP 204
}
