package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type EventTemplate struct {
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	Title      string            `json:"title" gorm:"not null" example:"KDEN FNO Template"`
	Positions  DefaultPositions  `json:"positions" gorm:"type:json" example:"[\"ZDV_APP\", \"ZDV_TWR\"]"`
	Facilities DefaultFacilities `json:"facilities" gorm:"type:json" example:"[\"ZDV\", \"ZAB\", \"ZLC\"]"`
	Fields     Fields            `json:"fields" gorm:"type:json" example:"[\"KDEN\", \"KBJC\", \"KAPA\"]"`
	Shifts     bool              `json:"shifts" gorm:"not null;default:false" example:"true"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type DefaultPositions []string

func (f *DefaultPositions) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f *DefaultPositions) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type DefaultFacilities []constants.FacilityID

func (f *DefaultFacilities) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f *DefaultFacilities) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (et *EventTemplate) Create() error {
	return database.DB.Create(et).Error
}

func (et *EventTemplate) Get() error {
	return database.DB.First(et).Error
}

func (et *EventTemplate) Update() error {
	return database.DB.Save(et).Error
}

func (et *EventTemplate) Delete() error {
	return database.DB.Delete(et).Error
}

func GetEventTemplatesFiltered(facility constants.FacilityID) ([]EventTemplate, error) {
	var ets []EventTemplate
	if err := database.DB.Where("facilities LIKE ?", "%"+facility+"%").Find(&ets).Error; err != nil {
		return nil, err
	}
	return ets, nil
}
