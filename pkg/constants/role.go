package constants

type RoleID string
type RoleManagementType uint64

const (
	UnknownRoleManagementType RoleManagementType = iota
	DivisionManagementManaged
	DivisionManaged
	FacilityManaged
)

// intentionally not capitalized to prevent direct usage outside this package
type role struct {
	Name                   string
	IsFacilitySpecificRole bool
	RoleManagementType     RoleManagementType
	// FacilityManagerRoles - Roles allowed to assign this role when RoleManagementType = FacilityManaged
	FacilityManagerRoles []RoleID
}

const (
	// ARTCC Roles
	AirTrafficManagerRole         RoleID = "ATM"
	DeputyAirTrafficManagerRole   RoleID = "DATM"
	TrainingAdministratorRole     RoleID = "TA"
	EventCoordinatorRole          RoleID = "EC"
	AssistantEventCoordinatorRole RoleID = "AEC"
	FacilityEngineerRole          RoleID = "FE"
	AssistantFacilityEngineerRole RoleID = "AFE"
	WebMasterRole                 RoleID = "WM"
	AssistantWebMasterRole        RoleID = "AWM"
	InstructorRole                RoleID = "INS"
	MentorRole                    RoleID = "MTR"

	// Division Roles
	DivisionStaffRole      RoleID = "DIVISION_STAFF"
	DivisionManagementRole RoleID = "DIVISION_MANAGEMENT"

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

var (
	// roleMap - Roles must exist in this map to be valid
	roleMap = map[RoleID]role{
		// Division Roles
		DivisionManagementRole: {
			Name:                   "Division Management",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManagementManaged,
		},
		DivisionStaffRole: {
			Name:                   "Division Staff",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManagementManaged,
		},
		DeveloperTeamRole: {
			Name:                   "Developer Team",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		AceTeamRole: {
			Name:                   "ACE Team",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		NTMSRole: {
			Name:                   "National Traffic Management Supervisor",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		NTMTRole: {
			Name:                   "National Traffic Management Team",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		SocialMediaTeam: {
			Name:                   "Social Media Team",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		TrainingContentTeam: {
			Name:                   "Training Content Team",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},
		AcademyMaterialEditor: {
			Name:                   "Academy Material Editor",
			IsFacilitySpecificRole: false,
			RoleManagementType:     DivisionManaged,
		},

		// FacilityID Roles
		AirTrafficManagerRole: {
			Name:                   "Air Traffic Manager",
			IsFacilitySpecificRole: true,
			RoleManagementType:     DivisionManaged,
		},
		DeputyAirTrafficManagerRole: {
			Name:                   "Deputy Air Traffic Manager",
			IsFacilitySpecificRole: true,
			RoleManagementType:     DivisionManaged,
		},
		TrainingAdministratorRole: {
			Name:                   "Training Administrator",
			IsFacilitySpecificRole: true,
			RoleManagementType:     DivisionManaged,
		},
		EventCoordinatorRole: {
			Name:                   "Event Coordinator",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole},
		},
		AssistantEventCoordinatorRole: {
			Name:                   "Assistant Event Coordinator",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole, EventCoordinatorRole},
		},
		FacilityEngineerRole: {
			Name:                   "Facility Engineer",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole},
		},
		AssistantFacilityEngineerRole: {
			Name:                   "Assistant Facility Engineer",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole, FacilityEngineerRole},
		},
		WebMasterRole: {
			Name:                   "Webmaster",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole},
		},
		AssistantWebMasterRole: {
			Name:                   "Assistant Webmaster",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole, WebMasterRole},
		},
		InstructorRole: {
			Name:                   "Instructor",
			IsFacilitySpecificRole: true,
			RoleManagementType:     DivisionManaged,
		},
		MentorRole: {
			Name:                   "Mentor",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles: []RoleID{
				AirTrafficManagerRole, DeputyAirTrafficManagerRole, TrainingAdministratorRole},
		},
		FacilityMaterialEditor: {
			Name:                   "FacilityID Material Editor",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles: []RoleID{
				AirTrafficManagerRole, DeputyAirTrafficManagerRole, TrainingAdministratorRole},
		},
		EmailUser: {
			Name:                   "Email User",
			IsFacilitySpecificRole: true,
			RoleManagementType:     FacilityManaged,
			FacilityManagerRoles:   []RoleID{AirTrafficManagerRole, DeputyAirTrafficManagerRole},
		},
	}
)

func (r RoleID) IsValidRole() bool {
	_, ok := roleMap[r]
	return ok
}

func (r RoleID) DisplayName() string {
	val, ok := roleMap[r]
	// Needed to prevent panic on invalid RoleID
	if ok {
		return val.Name
	}
	return ""
}

func (r RoleID) IsFacilitySpecific() bool {
	val, ok := roleMap[r]
	// Needed to prevent panic on invalid RoleID
	if ok {
		return val.IsFacilitySpecificRole
	}
	return false
}

func (r RoleID) ManagementType() RoleManagementType {
	val, ok := roleMap[r]
	if ok {
		return val.RoleManagementType
	}
	return UnknownRoleManagementType
}

func (r RoleID) FacilityManagerRoles() []RoleID {
	val, ok := roleMap[r]
	if ok {
		return val.FacilityManagerRoles
	}
	return []RoleID{}
}
