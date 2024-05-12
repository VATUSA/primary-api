package constants

type RoleID string
type GroupID string

type Role struct {
	Name         string
	Groups       []GroupID // Groups this role is a part of
	RolesCanAdd  []RoleID  // Roles this role can be added by
	GroupsCanAdd []GroupID // Groups this role can be added by
}

const (
	// ARTCC Roles
	AirTrafficManagerRole       RoleID = "ATM"
	DeputyAirTrafficManagerRole RoleID = "DATM"
	TrainingAdministratorRole   RoleID = "TA"
	EventCoordinatorRole        RoleID = "EC"
	AssistantEventCoordinator   RoleID = "AEC"
	FacilityEngineerRole        RoleID = "FE"
	AssistantFacilityEngineer   RoleID = "AFE"
	WebMasterRole               RoleID = "WM"
	AssistantWebMasterRole      RoleID = "AWM"
	InstructorRole              RoleID = "INS"
	MentorRole                  RoleID = "MTR"

	// Division Roles
	DivisionDirectorRole        RoleID = "USA1"
	AirTrafficServicesRole      RoleID = "USA2"
	TrainingServicesRole        RoleID = "USA3"
	SupportServicesRole         RoleID = "USA4"
	EventsManagerRole           RoleID = "USA5"
	TechnicalManagerRole        RoleID = "USA6"
	StaffDevelopmentManagerRole RoleID = "USA7"
	TrainingServicesManagerRole RoleID = "USA8"
	TrainingContentManagerRole  RoleID = "USA9"

	// Other Roles
	DeveloperTeamRole      RoleID = "DEV"
	AceTeamRole            RoleID = "ACE"
	NTMSRole               RoleID = "NTMS"
	NTMTRole               RoleID = "NTMT"
	SocialMediaTeam        RoleID = "SMT"
	TrainingContentTeam    RoleID = "TCT"
	AcademyMaterialEditor  RoleID = "CBT"
	FacilityMaterialEditor RoleID = "FACCBT"

	// Misc Roles
	EmailUser RoleID = "EMAIL"
)

const (
	// Division Groups
	DivisionManagement  GroupID = "div_mgmt"
	DivisionStaff       GroupID = "div_staff"
	DivisionDevelopment GroupID = "div_dev"
	DivisionEvents      GroupID = "div_events"
	DivisionTraining    GroupID = "div_training"

	// Facility Groups
	FacilityManagement  GroupID = "fac_mgmt"
	FacilityStaff       GroupID = "fac_staff"
	FacilityEvents      GroupID = "fac_events"
	FacilityTraining    GroupID = "fac_training"
	FacilityDevelopment GroupID = "fac_dev"
	FacilityEngineers   GroupID = "fac_eng"

	// Misc Groups
	TrafficManagement GroupID = "tmu"
)

var Roles = map[RoleID]Role{
	// ARTCC Roles
	AirTrafficManagerRole: {
		Name: "Air Traffic Manager",
		Groups: []GroupID{
			FacilityManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	DeputyAirTrafficManagerRole: {
		Name: "Deputy Air Traffic Manager",
		Groups: []GroupID{
			FacilityManagement,
		},
		RolesCanAdd: []RoleID{
			AirTrafficManagerRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			FacilityManagement,
		},
	},
	TrainingAdministratorRole: {
		Name:        "Training Administrator",
		Groups:      []GroupID{},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
			FacilityManagement,
		},
	},
	EventCoordinatorRole: {
		Name: "Event Coordinator",
		Groups: []GroupID{
			FacilityStaff,
			FacilityEvents,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionEvents,
			FacilityManagement,
		},
	},
	AssistantEventCoordinator: {
		Name: "Assistant Event Coordinator",
		Groups: []GroupID{
			FacilityEvents,
		},
		RolesCanAdd: []RoleID{
			EventCoordinatorRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionEvents,
			FacilityManagement,
		},
	},
	FacilityEngineerRole: {
		Name: "Facility Engineer",
		Groups: []GroupID{
			FacilityStaff,
			FacilityEngineers,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			FacilityManagement,
		},
	},
	AssistantFacilityEngineer: {
		Name: "Assistant Facility Engineer",
		Groups: []GroupID{
			FacilityEngineers,
		},
		RolesCanAdd: []RoleID{
			FacilityEngineerRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			FacilityManagement,
		},
	},
	WebMasterRole: {
		Name: "Webmaster",
		Groups: []GroupID{
			FacilityStaff,
			FacilityDevelopment,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			FacilityManagement,
		},
	},
	AssistantWebMasterRole: {
		Name: "Assistant Webmaster",
		Groups: []GroupID{
			FacilityDevelopment,
		},
		RolesCanAdd: []RoleID{
			WebMasterRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			FacilityManagement,
			FacilityDevelopment,
		},
	},
	InstructorRole: {
		Name: "Instructor",
		Groups: []GroupID{
			FacilityTraining,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionTraining,
			DivisionDevelopment,
		},
	},
	MentorRole: {
		Name: "Mentor",
		Groups: []GroupID{
			FacilityTraining,
		},
		RolesCanAdd: []RoleID{
			TrainingAdministratorRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionTraining,
			DivisionDevelopment,
			FacilityManagement,
		},
	},

	// Division Roles
	DivisionDirectorRole: {
		Name: "Division Director",
		Groups: []GroupID{
			DivisionManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	AirTrafficServicesRole: {
		Name: "Deputy Director Air Traffic Services",
		Groups: []GroupID{
			DivisionManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	TrainingServicesRole: {
		Name: "Deputy Director Training Services",
		Groups: []GroupID{
			DivisionManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	SupportServicesRole: {
		Name: "Deputy Director Support Services",
		Groups: []GroupID{
			DivisionManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	EventsManagerRole: {
		Name: "Events Manager",
		Groups: []GroupID{
			DivisionStaff,
			DivisionEvents,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	TechnicalManagerRole: {
		Name: "Technical Manager",
		Groups: []GroupID{
			DivisionStaff,
			DivisionDevelopment,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	StaffDevelopmentManagerRole: {
		Name: "Staff Development Manager",
		Groups: []GroupID{
			DivisionStaff,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	TrainingServicesManagerRole: {
		Name: "Training Services Manager",
		Groups: []GroupID{
			DivisionStaff,
			DivisionTraining,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
		},
	},
	TrainingContentManagerRole: {
		Name: "Training Content Manager",
		Groups: []GroupID{
			DivisionStaff,
			DivisionTraining,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
		},
	},

	// Other Roles
	DeveloperTeamRole: {
		Name: "Developer Team",
		Groups: []GroupID{
			DivisionDevelopment,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	AceTeamRole: {
		Name: "ACE Team",
		Groups: []GroupID{
			TrafficManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionEvents,
		},
	},
	NTMSRole: {
		Name: "National Traffic Management Supervisor",
		Groups: []GroupID{
			TrafficManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionEvents,
		},
	},
	NTMTRole: {
		Name: "National Training Management Supervisor",
		Groups: []GroupID{
			TrafficManagement,
		},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionEvents,
		},
	},
	SocialMediaTeam: {
		Name:        "Social Media Team",
		Groups:      []GroupID{},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
		},
	},
	TrainingContentTeam: {
		Name:        "Training Content Team",
		Groups:      []GroupID{},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
		},
	},
	AcademyMaterialEditor: {
		Name:        "Academy Material Editor (Academy)",
		Groups:      []GroupID{},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
		},
	},
	FacilityMaterialEditor: {
		Name:   "Academy Material Editor (Facility)",
		Groups: []GroupID{},
		RolesCanAdd: []RoleID{
			TrainingAdministratorRole,
		},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionDevelopment,
			DivisionTraining,
			FacilityManagement,
		},
	},

	// Misc Roles
	EmailUser: {
		Name:        "Email User",
		Groups:      []GroupID{},
		RolesCanAdd: []RoleID{},
		GroupsCanAdd: []GroupID{
			DivisionManagement,
			DivisionStaff,
			FacilityManagement,
			FacilityStaff,
		},
	},
}

func (r RoleID) IsValidRole() bool {
	_, ok := Roles[r]
	return ok
}

func (r RoleID) DisplayName() string {
	return Roles[r].Name
}

func (g GroupID) RolesInGroup() []RoleID {
	rolesInGroup := []RoleID{}
	for role, _ := range Roles {
		if role.InGroup(g) {
			rolesInGroup = append(rolesInGroup, role)
		}
	}
	return rolesInGroup

}

func (r RoleID) RolesCanAdd() []RoleID {
	if !r.IsValidRole() {
		return []RoleID{}
	}

	return Roles[r].RolesCanAdd
}

func (r RoleID) InGroup(group GroupID) bool {
	if !r.IsValidRole() {
		return false
	}
	for _, g := range Roles[r].Groups {
		if g == group {
			return true
		}
	}
	return false
}

// CanAddRole checks if a user with the given roles can add the target role
func CanAddRole(userRoles []RoleID, targetRole RoleID) bool {
	if !targetRole.IsValidRole() {
		return false
	}

	targettedRole := Roles[targetRole]
	for _, role := range userRoles {
		for _, r := range targettedRole.RolesCanAdd {
			if role == r {
				return true
			}
		}
		for _, group := range targettedRole.GroupsCanAdd {
			if role.InGroup(group) {
				return true
			}
		}
	}

	return false
}
