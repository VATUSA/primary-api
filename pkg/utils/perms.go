package utils

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database/models"
)

func IsVATUSAStaff(user *models.User) bool {
	for _, roster := range user.Roster {
		if roster.Facility == "ZHQ" {
			return true
		}
	}

	return false
}

func IsSeniorStaff(user *models.User) bool {
	for _, roster := range user.Roster {
		for _, roles := range roster.Roles {
			if roles.RoleID == constants.AirTrafficManagerRole || roles.RoleID == constants.DeputyAirTrafficManagerRole || roles.RoleID == constants.TrainingAdministratorRole {
				return true
			}
		}
	}

	return false
}

func IsFacilitySeniorStaff(user *models.User, facility constants.FacilityID) bool {
	for _, roster := range user.Roster {
		if roster.Facility == facility {
			for _, roles := range roster.Roles {
				if roles.RoleID == constants.AirTrafficManagerRole || roles.RoleID == constants.DeputyAirTrafficManagerRole || roles.RoleID == constants.TrainingAdministratorRole {
					return true
				}
			}
		}
	}

	return false
}

func IsFacilityStaff(user *models.User, facility constants.FacilityID) bool {
	for _, roster := range user.Roster {
		if roster.Facility == facility {
			for _, roles := range roster.Roles {
				if roles.RoleID == constants.FacilityEngineerRole || roles.RoleID == constants.WebMasterRole || roles.RoleID == constants.EventCoordinatorRole || IsFacilitySeniorStaff(user, facility) {
					return true
				}
			}
		}
	}

	return false
}

func IsFacilityEventsStaff(user *models.User, facility constants.FacilityID) bool {
	for _, roster := range user.Roster {
		if roster.Facility == facility {
			for _, roles := range roster.Roles {
				if roles.RoleID == constants.EventCoordinatorRole || roles.RoleID == constants.AssistantEventCoordinator {
					return true
				}
			}
		}
	}

	return false
}

func CanEditFacility(user *models.User, facility constants.FacilityID) bool {
	for _, roster := range user.Roster {
		if roster.Facility == facility {
			for _, roles := range roster.Roles {
				if roles.RoleID == constants.AirTrafficManagerRole || roles.RoleID == constants.DeputyAirTrafficManagerRole || roles.RoleID == constants.WebMasterRole {
					return true
				}
			}
		}
	}

	return false
}

func CanAddRole(user *models.User, roleId constants.RoleID, facilityId constants.FacilityID) bool {
	var userRoles []constants.RoleID
	for _, roster := range user.Roster {
		if roster.Facility == facilityId || roster.Facility == "ZHQ" {
			for _, role := range roster.Roles {
				userRoles = append(userRoles, role.RoleID)
			}
		}
	}

	return constants.CanAddRole(userRoles, roleId)
}

func IsInstructor(user *models.User, facility constants.FacilityID) bool {
	for _, roster := range user.Roster {
		if roster.Facility == facility {
			for _, roles := range roster.Roles {
				if roles.RoleID == constants.TrainingAdministratorRole || roles.RoleID == constants.InstructorRole {
					return user.ControllerRating == constants.InstructorRating || user.ControllerRating == constants.SeniorInstructorRating
				}
			}
		}
	}

	return false
}
