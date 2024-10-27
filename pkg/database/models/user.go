package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"strings"
	"time"
)

type User struct {
	CID                  uint                   `gorm:"primaryKey" json:"cid" example:"1293257"`
	FirstName            string                 `json:"first_name" example:"Raaj" gorm:"index:idx_first_name"`
	LastName             string                 `json:"last_name" example:"Patel" gorm:"index:idx_last_name"`
	PreferredName        string                 `json:"preferred_name" example:"Raaj" gorm:"index:idx_pref_name"`
	PrefNameEnabled      bool                   `json:"pref_name_enabled" example:"true"`
	Email                string                 `json:"email" example:"vatusa6@vatusa.net"`
	PreferredOIs         string                 `json:"preferred_ois" gorm:"column:preferred_ois" example:"RP"`
	PilotRating          constants.PilotRating  `json:"pilot_rating" example:"1"`
	ControllerRating     constants.ATCRating    `json:"controller_rating" example:"1"`
	DiscordID            string                 `json:"discord_id" example:"1234567890"`
	LastLogin            time.Time              `json:"last_login" example:"2021-01-01T00:00:00Z"`
	LastCertSync         time.Time              `json:"last_cert_sync" example:"2021-01-01T00:00:00Z"`
	Flags                UserFlag               `json:"flags" gorm:"foreignKey:CID"`
	RatingChanges        []RatingChange         `json:"-" gorm:"foreignKey:CID"`
	RosterRequest        []RosterRequest        `json:"-" gorm:"foreignKey:CID"`
	Roster               []Roster               `json:"-" gorm:"foreignKey:CID"`
	Notifications        []Notification         `json:"-" gorm:"foreignKey:CID"`
	Feedback             []Feedback             `json:"-" gorm:"foreignKey:ControllerCID"`
	ActionLogEntry       []ActionLogEntry       `json:"-" gorm:"foreignKey:CID"`
	DisciplinaryLogEntry []DisciplinaryLogEntry `json:"-" gorm:"foreignKey:CID"`
	CreatedAt            time.Time              `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt            time.Time              `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (un *User) Create() error {
	return database.DB.Create(un).Error
}

func (un *User) Update() error {
	return database.DB.Save(un).Error
}

func (un *User) Delete() error {
	return database.DB.Delete(un).Error
}

func (un *User) Get() error {
	if un.Email != "" {
		return database.DB.Where("email = ?", un.Email).First(un).Error
	}

	if un.DiscordID != "" {
		return database.DB.Where("discord_id = ?", un.DiscordID).First(un).Error
	}

	return database.DB.Preload("Roster").First(un).Error
}

func GetAllUsers() ([]User, error) {
	var users []User
	return users, database.DB.Find(&users).Error
}

func SearchUsersByName(query string) ([]User, error) {
	var users []User

	// Split the query into parts
	queryParts := strings.Fields(query)

	// Using LIKE condition for case-insensitive partial matching on both first name and last name
	for _, part := range queryParts {
		if err := database.DB.Where("lower(first_name) LIKE ?", "%"+strings.ToLower(part)+"%").
			Or("lower(last_name) LIKE ?", "%"+strings.ToLower(part)+"%").
			Or("lower(preferred_name) LIKE ?", "%"+strings.ToLower(part)+"%").
			Find(&users).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

func IsValidUser(cid uint) bool {
	var user User
	if err := database.DB.Where("cid = ?", cid).First(&user).Error; err != nil {
		return false
	}
	return true
}

func GetUserOIs(cid uint, facility constants.FacilityID) (string, error) {
	var user User
	if err := database.DB.Preload("Roster").First(&user, cid).Error; err != nil {
		return "", err
	}

	for _, roster := range user.Roster {
		if roster.Facility == facility {
			return roster.OIs, nil
		}
	}
	return user.PreferredOIs, nil
}
