package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/types"
	"gorm.io/gorm"
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

func (rr *RosterRequest) BeforeUpdate(tx *gorm.DB) error {
	oldRR := &RosterRequest{ID: rr.ID}
	if err := oldRR.Get(); err != nil {
		return err
	}
	if oldRR.Status == types.Pending && rr.Status == types.Accepted {
		roster := &Roster{
			CID:      rr.CID,
			Facility: rr.Facility,
			OIs:      "",
			Home:     false,
			Visiting: false,
			Status:   "Active",
		}

		if rr.RequestType == types.Visiting {
			roster.Visiting = true
		} else {
			roster.Home = true
		}

		if err := roster.Create(); err != nil {
			return err
		}

		// Transfers:
		// On accepting a transfer remove the user from their current facility
		if rr.RequestType == types.Transferring {
			rosters, err := GetRostersByCID(rr.CID)
			if err != nil {
				return err
			}

			for _, r := range rosters {
				if r.Facility != rr.Facility && r.Home {
					// if the user is an assistant add them to the new facility as a visitor
					for _, role := range r.Roles {
						if role.RoleID.IsAssistant() {
							roster := &Roster{
								CID:      r.CID,
								Facility: r.Facility,
								OIs:      r.OIs,
								Home:     false,
								Visiting: true,
								Status:   "Active",
							}

							if err := roster.Create(); err != nil {
								return err
							}
							break
						}
					}

					if err := r.Delete(); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rr *RosterRequest) Create() error {
	return database.DB.Create(rr).Error
}

func (rr *RosterRequest) Update() error {
	return database.DB.Updates(rr).Error
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

func GetFilteredRosterRequests(cid uint, reqType types.RequestType, dateAfter time.Time) ([]RosterRequest, error) {
	var rosterRequests []RosterRequest

	query := database.DB
	if cid != 0 {
		query = query.Where("cid = ?", cid)
	}
	if reqType != "" {
		query = query.Where("request_type = ?", reqType)
	}
	if !dateAfter.IsZero() {
		query = query.Where("created_at > ?", dateAfter)
	}

	return rosterRequests, query.Find(&rosterRequests).Error
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
