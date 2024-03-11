package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

// Expire Time can be the time of the session, or the time of the event

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	CID       uint      `json:"cid" example:"1293257"`
	Category  string    `json:"category" example:"Training"`
	Title     string    `json:"title" example:"Upcoming Training Session"`
	Body      string    `json:"body" example:"You have a training session coming up."`
	ExpireAt  time.Time `json:"expire_at" example:"2021-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (n *Notification) Create() error {
	return database.DB.Create(n).Error
}

func (n *Notification) Update() error {
	return database.DB.Save(n).Error
}

func (n *Notification) Delete() error {
	return database.DB.Delete(n).Error
}

func (n *Notification) Get() error {
	return database.DB.Where("id = ?", n.ID).First(n).Error
}

func GetAllNotifications() ([]Notification, error) {
	var notifications []Notification
	return notifications, database.DB.Find(&notifications).Error
}

func GetAllActiveNotificationsByCID(cid uint) ([]Notification, error) {
	var notifications []Notification
	return notifications, database.DB.Where("cid = ? AND expire_at > ?", cid, time.Now()).Find(&notifications).Error
}
