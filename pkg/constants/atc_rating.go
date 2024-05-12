package constants

type ATCRating int

// intentionally not capitalized to prevent direct usage outside this package
type atcRating struct {
	Short string
	Long  string
}

const (
	InactiveRating         ATCRating = -1
	SuspendedRating        ATCRating = 0
	ObserverRating         ATCRating = 1
	Student1Rating         ATCRating = 2
	Student2Rating         ATCRating = 3
	Student3Rating         ATCRating = 4
	ControllerRating       ATCRating = 5
	SeniorControllerRating ATCRating = 7
	InstructorRating       ATCRating = 8
	SeniorInstructorRating ATCRating = 10
	SupervisorRating       ATCRating = 11
	AdministratorRating    ATCRating = 12
)

var (
	atcRatingMap = map[ATCRating]atcRating{
		InactiveRating: {
			Short: "AFK",
			Long:  "Inactive",
		},
		SuspendedRating: {
			Short: "SUS",
			Long:  "Suspended",
		},
		ObserverRating: {
			Short: "OBS",
			Long:  "Observer",
		},
		Student1Rating: {
			Short: "S1",
			Long:  "Student 1",
		},
		Student2Rating: {
			Short: "S2",
			Long:  "Student 2",
		},
		Student3Rating: {
			Short: "S3",
			Long:  "Student 3",
		},
		ControllerRating: {
			Short: "C1",
			Long:  "Controller",
		},
		SeniorControllerRating: {
			Short: "C3",
			Long:  "Senior Controller",
		},
		InstructorRating: {
			Short: "I1",
			Long:  "Instructor",
		},
		SeniorInstructorRating: {
			Short: "I3",
			Long:  "Senior Instructor",
		},
		SupervisorRating: {
			Short: "SUP",
			Long:  "Supervisor",
		},
		AdministratorRating: {
			Short: "ADM",
			Long:  "Administrator",
		},
	}
)

func (r ATCRating) IsValidRating() bool {
	_, ok := atcRatingMap[r]
	return ok
}

func (r ATCRating) Int() int {
	return int(r)
}

func (r ATCRating) Short() string {
	val, ok := atcRatingMap[r]
	if ok {
		return val.Short
	}
	return "UNK"
}

func (r ATCRating) Long() string {
	val, ok := atcRatingMap[r]
	if ok {
		return val.Long
	}
	return "Unknown"
}
