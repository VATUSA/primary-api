package user

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/VATUSA/primary-api/pkg/database/types"
	vatsim_api "github.com/VATUSA/primary-api/pkg/vatsim/api"
	"time"
)

// TODO 50/50 Rule

type TransferEligibility struct {
	IsVATUSAMember           bool
	NeedS1OrRCE              bool
	NoTransfersWithin90Days  bool
	NoOtherTransfers         bool
	NoPromotionsWithin90Days bool
	NoStaffAtCurrentFacility bool
	NotInstructor            bool
	UsedTransferWaiver       bool
}

func getTransferEligible(user *models.User) (*TransferEligibility, error) {
	eligibility := &TransferEligibility{
		IsVATUSAMember: false,
	}

	location, err := vatsim_api.GetLocation(user.CID)
	if err != nil {
		return nil, err
	}

	if location.Division == "USA" {
		eligibility.IsVATUSAMember = true
	}

	// Needs RCE or S1 Exam
	// TODO: Implement

	// No Transfers Within 60 Days
	rosterReqs, err := models.GetFilteredRosterRequests(user.CID, types.Transferring, time.Now().AddDate(0, 0, -90))
	if err != nil {
		return nil, err
	}

	// No Other Transfers
	for _, rosterReq := range rosterReqs {
		if rosterReq.Status == types.Accepted {
			eligibility.NoTransfersWithin90Days = false
		}

		if rosterReq.Status == types.Pending {
			eligibility.NoOtherTransfers = false
		}
	}

	// No Promotions Within 90 Days
	ratingChanges, err := models.GetFilteredRatingChanges(user.CID, time.Now().Add(-90*24*time.Hour))
	if err != nil {
		return nil, err
	}

	for _, rc := range ratingChanges {
		if rc.NewRating == constants.SeniorControllerRating || rc.OldRating == constants.SeniorControllerRating {
			continue
		}
		if rc.NewRating == constants.InstructorRating || rc.OldRating == constants.InstructorRating {
			continue
		}
		if rc.NewRating == constants.SeniorInstructorRating || rc.OldRating == constants.SeniorInstructorRating {
			continue
		}
		if rc.NewRating == constants.SupervisorRating || rc.OldRating == constants.SupervisorRating {
			continue
		}
		if rc.NewRating == constants.AdministratorRating || rc.OldRating == constants.AdministratorRating {
			continue
		}

		if rc.NewRating > rc.OldRating {
			eligibility.NoPromotionsWithin90Days = false
			break
		}
	}

	// No Staff at Current Facility
	roles, err := models.GetAllUserRolesByCID(user.CID)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.RoleID.IsFacilityStaff() {
			eligibility.NoStaffAtCurrentFacility = false
			break
		}
	}

	// Not Instructor
	if user.ControllerRating == constants.InstructorRating || user.ControllerRating == constants.SeniorInstructorRating {
		eligibility.NotInstructor = false
	}

	// Used Transfer Waiver
	userFlags := &models.UserFlag{
		CID: user.CID,
	}
	if err := userFlags.Get(); err != nil {
		return nil, err
	}

	if userFlags.UsedTransferOverride {
		eligibility.UsedTransferWaiver = true
	}

	return eligibility, nil
}

type VisitingEligibility struct {
	HasS3                    bool
	NeedsRCE                 bool
	NoVisitsWithin60Days     bool
	NoOtherVisitingRequests  bool
	NoPromotionsWithin90Days bool
	Minimum50Hours           bool
	HasHomeFacility          bool
	VisitingAllowed          bool
}

func getVisitingEligible(user *models.User) (*VisitingEligibility, error) {
	eligibility := &VisitingEligibility{
		HasS3:                    true,
		NeedsRCE:                 true,
		NoVisitsWithin60Days:     true,
		NoOtherVisitingRequests:  true,
		NoPromotionsWithin90Days: true,
		Minimum50Hours:           true,
		HasHomeFacility:          true,
		VisitingAllowed:          true,
	}

	// Has S3
	if user.ControllerRating < constants.Student3Rating {
		eligibility.HasS3 = false
	}

	// Needs RCE
	// TODO: Implement

	// No Visits Within 60 Days
	rosterReqs, err := models.GetFilteredRosterRequests(user.CID, types.Visiting, time.Now().AddDate(0, 0, -60))
	if err != nil {
		return nil, err
	}

	for _, rosterReq := range rosterReqs {
		if rosterReq.Status == types.Pending {
			eligibility.NoOtherVisitingRequests = false
		}

		if rosterReq.Status == types.Accepted {
			eligibility.NoVisitsWithin60Days = false
		}
	}

	// No Promotions Within 90 Days
	ratingChanges, err := models.GetFilteredRatingChanges(user.CID, time.Now().Add(-90*24*time.Hour))
	if err != nil {
		return nil, err
	}

	for _, rc := range ratingChanges {
		if rc.NewRating == constants.SeniorControllerRating || rc.OldRating == constants.SeniorControllerRating {
			continue
		}
		if rc.NewRating == constants.InstructorRating || rc.OldRating == constants.InstructorRating {
			continue
		}
		if rc.NewRating == constants.SeniorInstructorRating || rc.OldRating == constants.SeniorInstructorRating {
			continue
		}
		if rc.NewRating == constants.SupervisorRating || rc.OldRating == constants.SupervisorRating {
			continue
		}
		if rc.NewRating == constants.AdministratorRating || rc.OldRating == constants.AdministratorRating {
			continue
		}

		if rc.NewRating > rc.OldRating {
			eligibility.NoPromotionsWithin90Days = false
			break
		}
	}

	// Minimum 50 Hours
	hours, err := vatsim_api.GetHours(user.CID)
	if err != nil {
		return nil, err
	}

	if user.ControllerRating == constants.Student1Rating && hours.S1 < 50 {
		eligibility.Minimum50Hours = false
	} else if user.ControllerRating == constants.Student2Rating && hours.S2 < 50 {
		eligibility.Minimum50Hours = false
	} else if user.ControllerRating == constants.Student3Rating && hours.S3 < 50 {
		eligibility.Minimum50Hours = false
	} else if user.ControllerRating == constants.ControllerRating && hours.C1 < 50 {
		eligibility.Minimum50Hours = false
	}

	// Visiting Allowed
	userFlags := &models.UserFlag{
		CID: user.CID,
	}
	if err := userFlags.Get(); err != nil {
		return nil, err
	}

	if userFlags.NoVisiting {
		eligibility.VisitingAllowed = false
	}

	return eligibility, nil
}
