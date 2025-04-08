package models

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database"
	"gorm.io/datatypes"
	"time"
)

type OTSRecord struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	StudentCID    uint           `json:"student_cid" gorm:"index"`
	InstructorCID uint           `json:"instructor_cid" gorm:"index"`
	Data          datatypes.JSON `json:"data" gorm:"type:jsonb"` // Stores form fields dynamically
	Notes         string         `json:"notes"`
	Result        bool           `json:"result"` // True = Passed, False = Failed
	CreatedAt     time.Time      `json:"created_at"`
}

func (otsRecord *OTSRecord) BeforeCreate() error {
	// check if student CID is valid
	if !IsValidUser(otsRecord.StudentCID) {
		return errors.New("invalid student CID")
	}

	return nil
}

func (otsRecord *OTSRecord) Create() error {
	return database.DB.Create(otsRecord).Error
}

func (otsRecord *OTSRecord) Update() error {
	return database.DB.Updates(otsRecord).Error
}

func (otsRecord *OTSRecord) Delete() error {
	return database.DB.Delete(otsRecord).Error
}

func (otsRecord *OTSRecord) Get() error {
	return database.DB.Where("id = ?", otsRecord.ID).First(otsRecord).Error
}

func GetFilteredOTSRecords(filter map[string]interface{}) ([]OTSRecord, error) {
	var records []OTSRecord
	err := database.DB.Where(filter).Find(&records).Error
	return records, err
}

type OTSTemplate struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name"`
	Template      datatypes.JSON `json:"template" gorm:"type:jsonb"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LastUpdatedBy uint           `json:"last_updated_by"`
}

func (otsTemplates *OTSTemplate) Create() error {
	return database.DB.Create(otsTemplates).Error
}

func (otsTemplates *OTSTemplate) Update() error {
	return database.DB.Updates(otsTemplates).Error
}

func (otsTemplates *OTSTemplate) Delete() error {
	return database.DB.Delete(otsTemplates).Error
}

func (otsTemplates *OTSTemplate) Get() error {
	return database.DB.Where("id = ?", otsTemplates.ID).First(otsTemplates).Error
}

func GetAllOTSTemplates() ([]OTSTemplate, error) {
	var templates []OTSTemplate
	err := database.DB.Find(&templates).Error
	return templates, err
}
