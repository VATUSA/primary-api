package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type TrainingNotes struct {
	ID            uint                 `json:"id" gorm:"primaryKey" example:"1"`
	StudentCID    uint                 `json:"student_cid" example:"1293257" gorm:"index"`
	InstructorCID uint                 `json:"instructor_cid" example:"1293257"`
	Facility      constants.FacilityID `json:"facility" example:"ZDV"`
	Position      string               `json:"position" example:"DEN_DEL"`
	Duration      time.Duration        `json:"duration" example:"1h30m"`
	Score         uint                 `json:"score" example:"100"`
	Notes         string               `json:"notes" example:"Great job!"`
	OTSRecordID   uint                 `json:"ots_record_id" example:"1"`
	SessionDate   time.Time            `json:"session_date" example:"2021-01-01T00:00:00Z"`
	CreatedAt     time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (tn *TrainingNotes) Create() error {
	return database.DB.Create(tn).Error
}

func (tn *TrainingNotes) Update() error {
	return database.DB.Save(tn).Error
}

func (tn *TrainingNotes) Delete() error {
	return database.DB.Delete(tn).Error
}

func (tn *TrainingNotes) Get() error {
	return database.DB.Where("id = ?", tn.ID).First(tn).Error
}

func GetFilteredTrainingNotes(filter map[string]interface{}) ([]TrainingNotes, error) {
	var notes []TrainingNotes
	err := database.DB.Where(filter).Find(&notes).Error
	return notes, err
}
