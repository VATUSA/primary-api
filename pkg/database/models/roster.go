package models

import (
	"errors"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type Roster struct {
	ID        uint                 `json:"id" gorm:"primaryKey" example:"1"`
	CID       uint                 `json:"cid" example:"1293257"`
	Facility  constants.FacilityID `json:"facility" example:"ZDV"`
	OIs       string               `json:"operating_initials" gorm:"column:ois" example:"RP"`
	Home      bool                 `json:"home" example:"true"`
	Visiting  bool                 `json:"visiting" example:"false"`
	Status    string               `json:"status" example:"Active"` // Active, LOA
	Roles     []UserRole           `json:"roles" gorm:"foreignKey:RosterID"`
	CreatedAt time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt *time.Time           `json:"deleted_at" example:"2021-01-01T00:00:00Z"` // Soft Deletes for logging
}

func (r *Roster) Create() error {
	// Check and see if user is already on the roster
	if err := database.DB.Where("c_id = ? AND facility = ?", r.CID, r.Facility).First(&Roster{}).Error; err == nil {
		return errors.New("user already exists on facility roster")
	}

	user := &User{CID: r.CID}
	if err := user.Get(); err != nil {
		return errors.New("user not found")
	}

	// See if preferred OIs are already taken
	if err := database.DB.Where("ois = ? AND facility = ?", user.PreferredOIs, r.Facility).First(&Roster{}).Error; err == nil {
		// OIs are taken so try first and last initial
		if err := database.DB.Where("ois = ? AND facility = ?", user.FirstName[:1]+user.LastName[:1], r.Facility).First(&Roster{}).Error; err == nil {
			// TODO - First and last initial are taken so just use first available OIs
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

func GetRosterByFacilityAndCID(facility constants.FacilityID, cid uint) (Roster, error) {
	var roster Roster
	return roster, database.DB.Where("facility = ? AND c_id = ?", facility, cid).First(&roster).Error
}

func GetRosters() ([]Roster, error) {
	var rosters []Roster
	return rosters, database.DB.Find(&rosters).Error
}

func GetRostersByCID(cid uint) ([]Roster, error) {
	var rosters []Roster
	return rosters, database.DB.Where("c_id = ?", cid).Preload("Roles").Find(&rosters).Error
}

func GetRostersByFacility(facility constants.FacilityID) ([]Roster, error) {
	var rosters []Roster
	return rosters, database.DB.Where("facility = ?", facility).Find(&rosters).Error
}

func GetRostersByFacilityAndType(facility constants.FacilityID, rosterType string) ([]Roster, error) {
	var rosters []Roster

	if rosterType == "home" {
		return rosters, database.DB.Where("facility = ? AND home = ?", facility, true).Find(&rosters).Error
	} else if rosterType == "visiting" {
		return rosters, database.DB.Where("facility = ? AND visiting = ?", facility, true).Find(&rosters).Error
	} else {
		return rosters, errors.New("invalid roster type")
	}
}
