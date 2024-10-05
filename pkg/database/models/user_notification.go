package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type UserNotification struct {
	CID            uint      `gorm:"primaryKey" json:"cid" example:"1293257"`
	DiscordEnabled bool      `json:"discord" example:"true"`
	EmailEnabled   bool      `json:"email" example:"true"`
	Events         bool      `json:"events" example:"true"`
	Training       bool      `json:"training" example:"true"`
	Feedback       bool      `json:"feedback" example:"true"`
	CreatedAt      time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (un *UserNotification) Create() error {
	return database.DB.Create(un).Error
}

func (un *UserNotification) Update() error {
	return database.DB.Save(un).Error
}

func (un *UserNotification) Delete() error {
	return database.DB.Delete(un).Error
}

func (un *UserNotification) Get() error {
	return database.DB.First(un, un.CID).Error
}
