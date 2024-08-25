package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"time"
)

type Feedback struct {
	ID            uint                 `json:"id" gorm:"primaryKey" example:"1"`
	PilotCID      uint                 `json:"pilot_cid" gorm:"column:pilot_cid" example:"1293257"`
	Pilot         User                 `json:"pilot" gorm:"foreignKey:PilotCID;references:CID"`
	Callsign      string               `json:"callsign" example:"DAL123"`
	ControllerCID uint                 `json:"controller_cid" gorm:"column:controller_cid" example:"1293257"`
	Controller    User                 `json:"controller" gorm:"foreignKey:ControllerCID;references:CID"`
	Position      string               `json:"position" example:"DEN_I_APP"`
	Facility      constants.FacilityID `json:"facility" example:"ZDV"`
	Rating        types.FeedbackRating `json:"rating" gorm:"type:enum('unsatisfactory', 'poor', 'fair', 'good', 'excellent');" example:"good"`
	Notes         string               `json:"notes" example:"Raaj was the best controller I've ever flown under."`
	Status        types.StatusType     `json:"status" gorm:"type:enum('pending', 'accepted', 'rejected');" example:"pending"`
	Comment       string               `json:"comment" example:"Great work Raaj!"`
	CreatedAt     time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt     time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (f *Feedback) Create() error {
	return database.DB.Create(f).Error
}

func (f *Feedback) Update() error {
	return database.DB.Save(f).Error
}

func (f *Feedback) Delete() error {
	return database.DB.Delete(f).Error
}

func (f *Feedback) Get() error {
	return database.DB.Where("id = ?", f.ID).First(f).Error
}

func GetAllFeedback(status types.StatusType) ([]Feedback, error) {
	var feedback []Feedback
	return feedback, database.DB.Find(&feedback).Error
}

func GetFeedbackByCID(cid uint, status types.StatusType) ([]Feedback, error) {
	var feedback []Feedback
	query := database.DB.Where("controller_cid = ?", cid)

	if status != types.All {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&feedback).Error
	return feedback, err
}

func GetFeedbackByFacility(facility constants.FacilityID, status types.StatusType) ([]Feedback, error) {
	if status != types.All {
		var feedback []Feedback
		return feedback, database.DB.Where("facility = ? AND status = ?", facility, status).Find(&feedback).Error
	}
	var feedback []Feedback
	return feedback, database.DB.Where("facility = ?", facility).Find(&feedback).Error
}

func GetFeedbackByFacilityAndCID(facility constants.FacilityID, controllerCID uint, status types.StatusType) ([]Feedback, error) {
	if status != types.All {
		var feedback []Feedback
		return feedback, database.DB.Where("facility = ? AND controller_cid = ? AND status = ?", facility, controllerCID, status).Find(&feedback).Error
	}
	var feedback []Feedback
	return feedback, database.DB.Where("facility = ? AND controller_cid = ?", facility, controllerCID).Find(&feedback).Error
}
