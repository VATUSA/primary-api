package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"time"
)

type RosterRequest struct {
	ID          uint                 `json:"id" gorm:"primaryKey" example:"1"`
	CID         uint                 `json:"cid" example:"1293257"`
	Facility    constants.FacilityID `json:"-" example:"ZDV"`
	RequestType types.RequestType    `json:"request_type" gorm:"type:enum('visiting', 'transferring');"`
	Status      types.StatusType     `json:"status" gorm:"type:enum('pending', 'accepted', 'rejected');"`
	Reason      string               `json:"reason" example:"I want to transfer to ZDV"`
	CreatedAt   time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt   time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (rr *RosterRequest) Create() error {
	return database.DB.Create(rr).Error
}

func (rr *RosterRequest) Update() error {
	return database.DB.Save(rr).Error
}

func (rr *RosterRequest) Delete() error {
	return database.DB.Delete(rr).Error
}

func (rr *RosterRequest) Get() error {
	return database.DB.Where("id = ?", rr.ID).First(rr).Error
}

func GetAllRosterRequests() ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Find(&rosterRequests).Error
}

func GetAllRosterRequestsByCID(cid uint) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Where("cid = ?", cid).Find(&rosterRequests).Error
}

func GetAllRosterRequestsByFacility(facility constants.FacilityID) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Where("facility = ?", facility).Find(&rosterRequests).Error
}

func GetRosterRequestsByType(facility constants.FacilityID, reqType types.RequestType) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Where("facility = ? AND request_type = ?", facility, reqType).Find(&rosterRequests).Error
}

func GetRosterRequestsByStatus(facility constants.FacilityID, status types.StatusType) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Where("facility = ? AND status = ?", facility, status).Find(&rosterRequests).Error
}

func GetRosterRequestsByTypeAndStatus(facility constants.FacilityID, reqType types.RequestType, status types.StatusType) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest
	return rosterRequests, database.DB.Where("facility = ? AND request_type = ? AND status = ?", facility, reqType, status).Find(&rosterRequests).Error
}
