package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type RatingChange struct {
	ID           uint                `json:"id" gorm:"primaryKey" example:"1"`
	CID          uint                `json:"cid" example:"1293257"`
	OldRating    constants.ATCRating `json:"old_rating" example:"1"`
	NewRating    constants.ATCRating `json:"new_rating" example:"2"`
	CreatedAt    time.Time           `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedByCID string              `json:"created_by_cid" example:"1293257"`
	UpdatedAt    time.Time           `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (rc *RatingChange) Create() error {
	return database.DB.Create(rc).Error
}

func (rc *RatingChange) Update() error {
	return database.DB.Save(rc).Error
}

func (rc *RatingChange) Delete() error {
	return database.DB.Delete(rc).Error
}

func (rc *RatingChange) Get() error {
	return database.DB.Where("id = ?", rc.ID).First(rc).Error
}

func GetAllRatingChanges() ([]RatingChange, error) {
	var ratingChanges []RatingChange
	return ratingChanges, database.DB.Find(&ratingChanges).Error
}

func GetAllRatingChangesByCID(cid uint) ([]RatingChange, error) {
	var ratingChanges []RatingChange
	return ratingChanges, database.DB.Where("cid = ?", cid).Find(&ratingChanges).Error
}
