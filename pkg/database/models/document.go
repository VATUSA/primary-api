package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"time"
)

type Document struct {
	ID          uint                   `json:"id" gorm:"primaryKey" example:"1"`
	Facility    constants.FacilityID   `json:"facility" example:"ZDV"`
	Name        string                 `json:"name" example:"DP001"`
	Description string                 `json:"description" example:"General Division Policy"`
	Category    types.DocumentCategory `gorm:"type:enum('general', 'training', 'information_technology', 'sops', 'loas', 'misc');" json:"category" example:"general"`
	URL         string                 `json:"url" example:"https://zdvartcc.org"`
	CreatedAt   time.Time              `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy   uint                   `json:"created_by" example:"1293257"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy   uint                   `json:"updated_by" example:"1293257"`
}

func (d *Document) Create() error {
	return database.DB.Create(d).Error
}

func (d *Document) Update() error {
	return database.DB.Save(d).Error
}

func (d *Document) Delete() error {
	return database.DB.Delete(d).Error
}

func (d *Document) Get() error {
	return database.DB.Where("id = ?", d.ID).First(d).Error
}

func GetAllDocuments() ([]Document, error) {
	var documents []Document
	return documents, database.DB.Find(&documents).Error
}

func GetAllDocumentsByCategory(category types.DocumentCategory) ([]Document, error) {
	var documents []Document
	return documents, database.DB.Where("category = ?", category).Find(&documents).Error
}

func GetAllDocumentsByFacility(facility constants.FacilityID) ([]Document, error) {
	var documents []Document
	return documents, database.DB.Where("facility = ?", facility).Find(&documents).Error
}

func GetAllDocumentsByFacilityAndCategory(facility constants.FacilityID, category types.DocumentCategory) ([]Document, error) {
	var documents []Document
	return documents, database.DB.Where("facility = ? AND category = ?", facility, category).Find(&documents).Error
}
