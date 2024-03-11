package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

// Possible ActionLogEntries types:
// - Changes to User Table
// - Changes to Roster Table (for given User)
// - Changes to User Role Table (for given User)
// - Changes to Rating Changes Table (for given User)

type ActionLogEntry struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	CID       uint      `json:"cid" example:"1293257"`
	Entry     string    `json:"entry" example:"Changed Preferred OIs to RP"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy string    `json:"created_by" example:"'1234567' or 'System'"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy string    `json:"updated_by" example:"'1234567' or 'System'"`
}

func (ale *ActionLogEntry) Create() error {
	return database.DB.Create(ale).Error
}

func (ale *ActionLogEntry) Update() error {
	return database.DB.Save(ale).Error
}

func (ale *ActionLogEntry) Delete() error {
	return database.DB.Delete(ale).Error
}

func (ale *ActionLogEntry) Get() error {
	return database.DB.Where("id = ?", ale.ID).First(ale).Error
}

func GetAllActionLogEntries() ([]ActionLogEntry, error) {
	var ale []ActionLogEntry
	return ale, database.DB.Find(&ale).Error
}

func GetAllActionLogEntriesByCID(cid uint) ([]ActionLogEntry, error) {
	var ale []ActionLogEntry
	return ale, database.DB.Where("cid = ?", cid).Find(&ale).Error
}
