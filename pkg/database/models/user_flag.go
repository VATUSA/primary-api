package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type UserFlag struct {
	CID                      uint      `json:"cid" example:"1293257"`
	NoStaffRole              bool      `json:"no_staff_role" example:"false"`
	NoStaffLogEntryID        uint      `json:"no_staff_log_entry_id" example:"1"`
	NoVisiting               bool      `json:"no_visiting" example:"false"`
	NoVisitingLogEntryID     uint      `json:"no_visiting_log_entry_id" example:"1"`
	NoTransferring           bool      `json:"no_transferring" example:"false"`
	NoTransferringLogEntryID uint      `json:"no_transferring_log_entry_id" example:"1"`
	NoTraining               bool      `json:"no_training" example:"false"`
	NoTrainingLogEntryID     uint      `json:"no_training_log_entry_id" example:"1"`
	CreatedAt                time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt                time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (f *UserFlag) Create() error {
	return database.DB.Create(f).Error
}

func (f *UserFlag) Update() error {
	return database.DB.Save(f).Error
}

func (f *UserFlag) Delete() error {
	return database.DB.Delete(f).Error
}

func (f *UserFlag) Get() error {
	return database.DB.Where("c_id = ?", f.CID).First(f).Error
}

func GetAllFlags() ([]UserFlag, error) {
	var flags []UserFlag
	return flags, database.DB.Find(&flags).Error
}
