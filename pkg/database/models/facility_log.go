package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

// Possible FacilityLogEntry types:
// - Changes to Facility Table
// - Changes to Faq Table (for given Facility)
// - Changes to Document Table (for given Facility)
// - Changes to Role Table (for given Facility)
// - Changes to gSuite Email (for given Facility)

type FacilityLogEntry struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Facility  string    `json:"facility" example:"ZDV"`
	Entry     string    `json:"entry" example:"Change URL from 'denartcc.org' to 'zdvartcc.org'"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy string    `json:"created_by" example:"'1234567' or 'System'"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy string    `json:"updated_by" example:"'1234567' or 'System'"`
}

func (fle *FacilityLogEntry) Create() error {
	return database.DB.Create(fle).Error
}

func (fle *FacilityLogEntry) Update() error {
	return database.DB.Save(fle).Error
}

func (fle *FacilityLogEntry) Delete() error {
	return database.DB.Delete(fle).Error
}

func (fle *FacilityLogEntry) Get() error {
	return database.DB.Where("id = ?", fle.ID).First(fle).Error
}

func GetAllFacilityLogEntries() ([]FacilityLogEntry, error) {
	var fle []FacilityLogEntry
	return fle, database.DB.Find(&fle).Error
}

func GetAllFacilityLogEntriesByFacility(facility string) ([]FacilityLogEntry, error) {
	var fle []FacilityLogEntry
	return fle, database.DB.Where("facility = ?", facility).Find(&fle).Error
}
