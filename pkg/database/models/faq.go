package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type FAQ struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Facility  string    `json:"facility" example:"ZDV"`
	Question  string    `json:"question" example:"Why shouldn't I join ZDV?'"`
	Answer    string    `json:"answer" example:"There are no reasons not to join ZDV!"`
	Category  string    `gorm:"type:enum('membership', 'training', 'technology', 'misc');" json:"category" example:"membership"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	CreatedBy uint      `json:"created_by" example:"1293257"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	UpdatedBy uint      `json:"updated_by" example:"1293257"`
}

func (f *FAQ) Create() error {
	return database.DB.Create(f).Error
}

func (f *FAQ) Update() error {
	return database.DB.Save(f).Error
}

func (f *FAQ) Delete() error {
	return database.DB.Delete(f).Error
}

func (f *FAQ) Get() error {
	return database.DB.Where("id = ?", f.ID).First(f).Error
}

func GetAllFAQ() ([]FAQ, error) {
	var faq []FAQ
	return faq, database.DB.Find(&faq).Error
}

func GetAllFAQByCategory(category string) ([]FAQ, error) {
	var faq []FAQ
	return faq, database.DB.Where("category = ?", category).Find(&faq).Error
}

func GetAllFAQByFacility(facility string) ([]FAQ, error) {
	var faq []FAQ
	return faq, database.DB.Where("facility = ?", facility).Find(&faq).Error
}
