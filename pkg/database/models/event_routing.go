package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type EventRouting struct {
	ID      uint `json:"id" gorm:"primaryKey" example:"1"`
	EventID uint `json:"event_id" gorm:"not null" example:"1"`

	Origin      string `json:"origin" example:"KLAX"`
	Destination string `json:"destination" example:"KDEN"`
	Routing     string `json:"routing" example:"DCT"`
	Notes       string `json:"notes" example:"JETS ONLY"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (er *EventRouting) Create() error {
	return database.DB.Create(er).Error
}

func (er *EventRouting) Get() error {
	return database.DB.First(er).Error
}

func (er *EventRouting) Update() error {
	return database.DB.Save(er).Error
}

func (er *EventRouting) Delete() error {
	return database.DB.Delete(er).Error
}

func GetEventRoutingFiltered(eventID uint) ([]EventRouting, error) {
	var routing []EventRouting
	query := database.DB
	if eventID != 0 {
		query = query.Where("event_id = ?", eventID)
	}
	return routing, query.Find(&routing).Error
}
