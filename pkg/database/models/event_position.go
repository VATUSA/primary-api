package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type EventPosition struct {
	ID      uint `json:"id" gorm:"primaryKey" example:"1"`
	EventID uint `json:"event_id" gorm:"not null" example:"1"`

	Facility          constants.FacilityID `json:"facility" example:"ZDV"`
	Position          string               `json:"position" gorm:"not null" example:"ZDV_APP"`
	Assignee          uint                 `json:"assignee" example:"1293257"`
	Shifts            bool                 `json:"shifts" gorm:"not null;default:false" example:"true"`
	SecondaryAssignee uint                 `json:"secondary_assignee" example:"1293257"`

	Signups []EventSignup `json:"signups" gorm:"foreignKey:PositionID"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (ep *EventPosition) Create() error {
	return database.DB.Create(ep).Error
}

func (ep *EventPosition) Get() error {
	return database.DB.Preload("Signups").First(ep).Error
}

func (ep *EventPosition) Update() error {
	return database.DB.Save(ep).Error
}

func (ep *EventPosition) Delete() error {
	return database.DB.Delete(ep).Error
}

func GetEventPositionsFiltered(eventId uint, facilityID constants.FacilityID) ([]EventPosition, error) {
	var positions []EventPosition
	query := database.DB.Where("event_id = ?", eventId)
	if facilityID != "" {
		query = query.Where("facility = ?", facilityID)
	}

	return positions, query.Find(&positions).Error
}
