package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type TrainingNotes struct {
	ID            uint                 `json:"id" gorm:"primaryKey" example:"1"`
	StudentCID    uint                 `json:"student_cid" example:"1293257" gorm:"index"`
	Student       *User                `json:"-" gorm:"foreignKey:StudentCID"`
	InstructorCID uint                 `json:"instructor_cid" example:"1293257"`
	Instructor    *User                `json:"-" gorm:"foreignKey:InstructorCID"`
	Facility      constants.FacilityID `json:"facility" example:"ZDV"`
	Position      string               `json:"position" example:"DEN_DEL"`
	Duration      time.Duration        `json:"duration" example:"60"`
	Notes         string               `json:"notes" example:"Great job!"`
	RatingExamID  uint                 `json:"rating_exam_id" example:"1"`
	RatingExam    *RatingExamRecord    `json:"-" gorm:"foreignKey:RatingExamID"`
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
	err := database.DB.Preload("Student").Preload("Instructor").Preload("RatingExam").Where(filter).Find(&notes).Error
	return notes, err
}
