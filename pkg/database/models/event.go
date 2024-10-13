package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type Event struct {
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	Title       string     `json:"title" gorm:"not null" example:"ZDV FNO"`
	Description string     `json:"description" gorm:"not null" example:"Join us for a fun night of flying in and out of Denver!"`
	BannerURL   string     `json:"banner_url" example:"https://zdvartcc.org/banner.jpg"`
	StartDate   time.Time  `json:"start_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	EndDate     time.Time  `json:"end_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	Fields      Fields     `json:"fields" gorm:"type:json" example:"[\"KDEN\", \"KBJC\", \"KAPA\"]"`
	Facilities  Facilities `json:"facilities" gorm:"type:json" example:"[\"ZDV\", \"ZAB\", \"ZLC\"]"`

	Positions []EventPosition `json:"positions" gorm:"foreignKey:EventID"`
	Routing   []EventRouting  `json:"routing" gorm:"foreignKey:EventID"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type Fields []string

func (f *Fields) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f *Fields) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type Facilities []constants.FacilityID

func (f *Facilities) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), f)
}

func (f *Facilities) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (e *Event) Create() error {
	return database.DB.Create(e).Error
}

func (e *Event) Get() error {
	return database.DB.Preload("Positions").Preload("Routing").First(e).Error
}

func (e *Event) Update() error {
	return database.DB.Save(e).Error
}

func (e *Event) Delete() error {
	return database.DB.Delete(e).Error
}

func GetEventsFiltered(facilityID constants.FacilityID, afterDate time.Time) ([]Event, error) {
	var events []Event
	query := database.DB
	if !afterDate.IsZero() {
		query = query.Where("start_date > ?", afterDate)
	}
	if facilityID != "" {
		// FIXME: idk if this query works
		query = query.Where("facilities @> ?", []constants.FacilityID{facilityID})
	}

	return events, query.Find(&events).Error
}