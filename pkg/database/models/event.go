package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type Event struct {
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	Title       string    `json:"title" gorm:"not null" example:"ZDV FNO"`
	Description string    `json:"description" gorm:"not null" example:"Join us for a fun night of flying in and out of Denver!"`
	BannerURL   string    `json:"banner_url" example:"https://zdvartcc.org/banner.jpg"`
	StartDate   time.Time `json:"start_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`
	EndDate     time.Time `json:"end_date" gorm:"not null" example:"2021-01-01T00:00:00Z"`

	Fields     []string               `json:"fields" gorm:"serializer:json" example:"[\"KDEN\", \"KBJC\", \"KAPA\"]"`
	Facilities []constants.FacilityID `json:"facilities" gorm:"serializer:json" example:"[\"ZDV\", \"ZAB\", \"ZLC\"]"`

	Positions []EventPosition `json:"positions" gorm:"foreignKey:EventID"`
	Routing   []EventRouting  `json:"routing" gorm:"foreignKey:EventID"`

	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (e *Event) Create() error {
	return database.DB.Create(e).Error
}

func (e *Event) Get() error {
	return database.DB.Preload("Positions.Signups").Preload("Routing").First(e).Error
}

func (e *Event) Update() error {
	return database.DB.Save(e).Error
}

func (e *Event) Delete() error {
	return database.DB.Delete(e).Error
}

func GetEventsFiltered(page, pageSize int, facilityID constants.FacilityID, afterDate time.Time) ([]Event, error) {
	var events []Event

	query := database.DB
	if !afterDate.IsZero() {
		query = query.Where("end_date > ?", afterDate)
	}
	if facilityID != "" {
		// FIXME: idk if this query works
		query = query.Where("facilities @> ?", []constants.FacilityID{facilityID})
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	return events, query.Find(&events).Error
}
