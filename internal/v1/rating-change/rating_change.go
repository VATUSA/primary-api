package rating_change

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	OldRating    uint   `json:"old_rating" example:"1" validate:"required"`
	NewRating    uint   `json:"new_rating" example:"2" validate:"required"`
	CreatedByCID string `json:"created_by_cid" example:"1293257" validate:"required"`
}

func (req *Request) Validate() error {
	return validator.New().Struct(req)
}

func (req *Request) Bind(r *http.Request) error {
	return nil
}

type Response struct {
	*models.RatingChange
}

func NewRatingChangeResponse(rc *models.RatingChange) *Response {
	return &Response{RatingChange: rc}
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	if res.RatingChange == nil {
		return errors.New("rating change not found")
	}
	return nil
}

func NewRatingChangeListResponse(rc []models.RatingChange) []render.Renderer {
	list := []render.Renderer{}
	for _, d := range rc {
		list = append(list, NewRatingChangeResponse(&d))
	}
	return list
}

// CreateRatingChange godoc
// @Summary Create a new rating change
// @Description Create a new rating change
// @Tags rating-change
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param rating_change body Request true "Rating Change"
// @Success 201 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/rating-change [post]
func CreateRatingChange(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	user := utils.GetUserCtx(r)

	if !models.IsValidUser(user.CID) {
		utils.Render(w, r, utils.ErrInvalidCID)
		return
	}

	rc := &models.RatingChange{
		CID:          user.CID,
		OldRating:    data.OldRating,
		NewRating:    data.NewRating,
		CreatedByCID: data.CreatedByCID,
	}

	if err := rc.Create(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, NewRatingChangeResponse(rc))
}

// ListRatingChanges godoc
// @Summary List rating changes
// @Description List rating changes
// @Tags rating-change
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Success 200 {object} []Response
// @Failure 422 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/rating-change [get]
func ListRatingChanges(w http.ResponseWriter, r *http.Request) {
	rc, err := models.GetAllRatingChangesByCID(utils.GetUserCtx(r).CID)
	if err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewRatingChangeListResponse(rc)); err != nil {
		utils.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateRatingChange godoc
// @Summary Update a rating change
// @Description Update a rating change
// @Tags rating-change
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param id path int true "Rating Change ID"
// @Param rating_change body Request true "Rating Change"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/rating-change/{id} [put]
func UpdateRatingChange(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if err := data.Validate(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	rc := utils.GetRatingChangeCtx(r)

	rc.OldRating = data.OldRating
	rc.NewRating = data.NewRating
	rc.CreatedByCID = data.CreatedByCID

	if err := rc.Update(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	utils.Render(w, r, NewRatingChangeResponse(rc))
}

// PatchRatingChange godoc
// @Summary Patch a rating change
// @Description Patch a rating change
// @Tags rating-change
// @Accept  json
// @Produce  json
// @Param cid path int true "CID"
// @Param id path int true "Rating Change ID"
// @Param rating_change body Request true "Rating Change"
// @Success 200 {object} Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/rating-change/{id} [patch]
func PatchRatingChange(w http.ResponseWriter, r *http.Request) {
	rc := utils.GetRatingChangeCtx(r)
	data := &Request{}
	if err := data.Bind(r); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	if data.OldRating != 0 {
		rc.OldRating = data.OldRating
	}
	if data.NewRating != 0 {
		rc.NewRating = data.NewRating
	}
	if data.CreatedByCID != "" {
		rc.CreatedByCID = data.CreatedByCID
	}

	if err := rc.Update(); err != nil {
		utils.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	utils.Render(w, r, NewRatingChangeResponse(rc))
}

// DeleteRatingChange godoc
// @Summary Delete a rating change
// @Description Delete a rating change
// @Tags rating-change
// @Accept  json
// @Produce  json
// @Param id path int true "Rating Change ID"
// @Success 204
// @Failure 400 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /user/{cid}/rating-change/{id} [delete]
func DeleteRatingChange(w http.ResponseWriter, r *http.Request) {
	rc := utils.GetRatingChangeCtx(r)
	if err := rc.Delete(); err != nil {
		utils.Render(w, r, utils.ErrInternalServer)
		return
	}

	render.Status(r, http.StatusNoContent)
}
