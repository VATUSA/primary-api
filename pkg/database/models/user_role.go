package models

import (
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"time"
)

type UserRole struct {
	ID         uint                 `json:"id" gorm:"primaryKey" example:"1"`
	CID        uint                 `json:"cid" example:"1293257"`
	RoleID     constants.RoleID     `json:"role" gorm:"type:varchar(10)" example:"ATM"`
	FacilityID constants.FacilityID `json:"facility_id" example:"ZDV"`
	CreatedAt  time.Time            `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt  time.Time            `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

func (ur *UserRole) Create() error {
	return database.DB.Create(ur).Error
}

func (ur *UserRole) Update() error {
	return database.DB.Save(ur).Error
}

func (ur *UserRole) Delete() error {
	return database.DB.Delete(ur).Error
}

func (ur *UserRole) Get() error {
	return database.DB.Where("id = ?", ur.ID).First(ur).Error
}

func GetAllUserRoles() ([]UserRole, error) {
	var userRoles []UserRole
	return userRoles, database.DB.Find(&userRoles).Error
}

func GetAllUserRolesByCID(cid uint) ([]UserRole, error) {
	var userRoles []UserRole
	return userRoles, database.DB.Where("cid = ?", cid).Find(&userRoles).Error
}

func GetAllUserRolesByRoleID(roleID string) ([]UserRole, error) {
	var userRoles []UserRole
	return userRoles, database.DB.Where("role_id = ?", roleID).Find(&userRoles).Error
}

func GetAllUserRolesByFacilityID(facilityID string) ([]UserRole, error) {
	var userRoles []UserRole
	return userRoles, database.DB.Where("facility_id = ?", facilityID).Find(&userRoles).Error
}

func CanModifyRole(user *User, role constants.RoleID, facilityId constants.FacilityID) bool {
	managementType := role.ManagementType()
	if managementType == constants.DivisionManagementManaged {
		return user.HasRoleAnyFacility(constants.DivisionManagementRole)
	} else if managementType == constants.DivisionManaged {
		return user.HasRoleListAnyFacility([]constants.RoleID{
			constants.DivisionManagementRole,
			constants.DivisionStaffRole,
		})
	} else if managementType == constants.FacilityManaged {
		return user.HasRoleListAnyFacility([]constants.RoleID{
			constants.DivisionManagementRole,
			constants.DivisionStaffRole,
		}) || user.HasRoleListAtFacility(role.FacilityManagerRoles(), facilityId)
	}

	return false
}
