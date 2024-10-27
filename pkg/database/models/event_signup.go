package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type EventSignup struct {
	ID      uint `json:"id" gorm:"primaryKey" example:"1"`
	EventID uint `json:"event_id" gorm:"not null" example:"1"`

	PositionID uint   `json:"position_id" gorm:"not null" example:"1"`
	CID        uint   `json:"cid" gorm:"not null" example:"1293257"`
	Name       string `json:"name" gorm:"-"`
	Shift      uint   `json:"shift" gorm:"not null;default:1" example:"1"` // 1 = Primary, 2 = Secondary

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (es *EventSignup) Create() error {
	return database.DB.Create(es).Error
}

func (es *EventSignup) Get() error {
	return database.DB.First(es).Error
}

func (es *EventSignup) Update() error {
	return database.DB.Save(es).Error
}

func (es *EventSignup) Delete() error {
	return database.DB.Delete(es).Error
}

func GetEventSignupFiltered(eventID uint) ([]EventSignup, error) {
	var signups []EventSignup
	query := database.DB
	if eventID != 0 {
		query = query.Where("event_id = ?", eventID)
	}
	return signups, query.Find(&signups).Error
}
