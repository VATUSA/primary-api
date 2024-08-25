package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

// Possible DisciplinaryLogEntry types:
// - Adding any flags
// - Changes to rating SUSPENDED

type DisciplinaryLogEntry struct {
	ID         uint      `json:"id" gorm:"primaryKey" example:"1"`
	CID        uint      `json:"cid" example:"1293257"`
	Entry      string    `json:"entry" example:"Changed Preferred OIs to RP"`
	VATUSAOnly bool      `json:"vatusa_only" example:"true"`
	CreatedAt  time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy  string    `json:"created_by" example:"'1234567' or 'System'"`
	UpdatedAt  time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy  string    `json:"updated_by" example:"'1234567' or 'System'"`
}

func (dle *DisciplinaryLogEntry) Create() error {
	return database.DB.Create(dle).Error
}

func (dle *DisciplinaryLogEntry) Update() error {
	return database.DB.Save(dle).Error
}

func (dle *DisciplinaryLogEntry) Delete() error {
	return database.DB.Delete(dle).Error
}

func (dle *DisciplinaryLogEntry) Get() error {
	return database.DB.Where("id = ?", dle.ID).First(dle).Error
}

func GetAllDisciplinaryLogEntries(VATUSAOnly bool) ([]DisciplinaryLogEntry, error) {
	var dle []DisciplinaryLogEntry
	return dle, database.DB.Where("vatusa_only = ?", VATUSAOnly).Find(&dle).Error
}

func GetAllDisciplinaryLogEntriesByCID(cid uint, VATUSAOnly bool) ([]DisciplinaryLogEntry, error) {
	if VATUSAOnly {
		var dle []DisciplinaryLogEntry
		return dle, database.DB.Where("cid = ?", cid).Find(&dle).Error
	}

	var dle []DisciplinaryLogEntry
	return dle, database.DB.Where("cid = ? AND vatusa_only = ?", cid, VATUSAOnly).Find(&dle).Error
}
