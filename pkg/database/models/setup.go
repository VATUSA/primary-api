package models

import (
	"github.com/VATUSA/primary-api/pkg/database"
	"log"
)

func AutoMigrate() {
	err := database.DB.AutoMigrate(
		&Facility{},
		&User{},
		&ActionLogEntry{},
		&DisciplinaryLogEntry{},
		&Document{},
		&FacilityLogEntry{},
		&FAQ{},
		&Feedback{},
		&News{},
		&Notification{},
		&RatingChange{},
		&Roster{},
		&RosterRequest{},
		&UserNotification{},
		&UserFlag{},
		&UserRole{},
	)
	if err != nil {
		log.Fatal("[Database] Migration Error:", err)
	}
}

func DropTables() {
	err := database.DB.Migrator().DropTable(
		&Facility{},
		&User{},
		&ActionLogEntry{},
		&DisciplinaryLogEntry{},
		&Document{},
		&FacilityLogEntry{},
		&FAQ{},
		&Feedback{},
		&News{},
		&Notification{},
		&RatingChange{},
		&Roster{},
		&RosterRequest{},
		&UserNotification{},
		&UserFlag{},
		&UserRole{},
	)
	if err != nil {
		log.Fatal("[Database] Drop Table Error:", err)
	}
}
