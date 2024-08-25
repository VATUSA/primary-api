package models

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type Facility struct {
	ID               constants.FacilityID `json:"id" gorm:"size:3;primaryKey" example:"ZDV"`
	Name             string               `json:"name" example:"Denver ARTCC"`
	URL              string               `json:"url" example:"https://zdvartcc.org"`
	APIKey           string               `json:"api_key" example:"1234567890"`
	FacilityLogEntry []FacilityLogEntry   `json:"-" gorm:"foreignKey:Facility"`
	FAQ              []FAQ                `json:"-" gorm:"foreignKey:Facility"`
	Document         []Document           `json:"-" gorm:"foreignKey:Facility"`
	CreatedAt        time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt        time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (f *Facility) Create() error {
	return database.DB.Create(f).Error
}

func (f *Facility) Update() error {
	return database.DB.Save(f).Error
}

func (f *Facility) Delete() error {
	return database.DB.Delete(f).Error
}

func (f *Facility) Get() error {
	return database.DB.Where("id = ?", f.ID).First(f).Error
}

func GetFacilityByAPIKey(key string) (*Facility, error) {
	var f Facility
	if err := database.DB.Where("api_key = ?", key).First(&f).Error; err != nil {
		return nil, err
	}

	return &f, nil
}

func IsValidFacility(id constants.FacilityID) bool {
	var f Facility
	return database.DB.Where("id = ?", id).First(&f).Error == nil
}

func GetAllFacilities() ([]Facility, error) {
	var facilities []Facility
	return facilities, database.DB.Find(&facilities).Error
}

func (f *Facility) StripSensitive() {
	f.APIKey = ""
}

func GenerateApiKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
