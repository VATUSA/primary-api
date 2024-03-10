package models

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/database"
	"gorm.io/gorm"
	"time"
)

type Roster struct {
	ID        uint       `json:"id" gorm:"primaryKey" example:"1"`
	CID       uint       `json:"cid" example:"1293257"`
	Facility  string     `json:"facility" example:"ZDV"`
	OIs       string     `json:"operating_initials" example:"RP"`
	Home      bool       `json:"home" example:"true"`
	Visiting  bool       `json:"visiting" example:"false"`
	Status    string     `json:"status" example:"Active"` // Active, LOA
	Roles     []UserRole `json:"roles" gorm:"foreignKey:ID"`
	CreatedAt time.Time  `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt time.Time  `json:"deleted_at" example:"2021-01-01T00:00:00Z"` // Soft Deletes for logging
}

func (r *Roster) Create() error {
	// Check and see if user is already on the roster\
	if err := database.DB.Where("cid = ? AND facility = ?", r.CID, r.Facility).First(&User{}).Error; err == nil {
		return errors.New("user already exists on facility roster")
	}

	user := &User{CID: r.CID}
	if err := user.Get(); err != nil {
		return errors.New("user not found")
	}

	// See if preferred OIs are already taken
	if err := database.DB.Where("ois = ? AND facility = ?", user.PreferredOIs, r.Facility).First(&User{}).Error; err == nil {
		// OIs are taken so try first and last initial
		if err := database.DB.Where("ois = ? AND facility = ?", user.FirstName[:1]+user.LastName[:1], r.Facility).First(&User{}).Error; err == nil {
			// First and last initial are taken so just use first available OIs
			return database.DB.Create(r).Error
		}
		r.OIs = user.FirstName[:1] + user.LastName[:1]
		return database.DB.Create(r).Error
	}

	r.OIs = user.PreferredOIs
	return database.DB.Create(r).Error
}

func (r *Roster) Update() error {
	return database.DB.Save(r).Error
}

func (r *Roster) Delete() error {
	return database.DB.Delete(r).Error
}

func (r *Roster) Get() error {
	return database.DB.Where("id = ?", r.ID).First(r).Error
}

func GetAllRosters() ([]Roster, error) {
	var rosters []Roster
	return rosters, database.DB.Find(&rosters).Error
}

func GetAllRostersByCID(db *gorm.DB, cid uint) ([]Roster, error) {
	var rosters []Roster
	return rosters, db.Where("cid = ?", cid).Find(&rosters).Error
}

func GetAllRostersByFacility(db *gorm.DB, facility string) ([]Roster, error) {
	var rosters []Roster
	return rosters, db.Where("facility = ?", facility).Find(&rosters).Error
}
