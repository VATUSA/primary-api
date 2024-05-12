package constants

type PilotRating int

type pilotRating struct {
	Short string
	Long  string
}

const (
	NoRating               PilotRating = 0
	PrivatePilotRating     PilotRating = 1
	InstrumentRating       PilotRating = 3
	CommercialMultiRating  PilotRating = 7
	AirTransportRating     PilotRating = 15
	FlightInstructorRating PilotRating = 31
	FlightExaminerRating   PilotRating = 63
)

var (
	pilotRatingMap = map[PilotRating]pilotRating{
		NoRating: {
			Short: "P0",
			Long:  "No Pilot Rating",
		},
		PrivatePilotRating: {
			Short: "PPL",
			Long:  "Private Pilot License",
		},
		InstrumentRating: {
			Short: "IR",
			Long:  "Instrument Rating",
		},
		CommercialMultiRating: {
			Short: "CMEL",
			Long:  "Commercial Multi-Engine License",
		},
		AirTransportRating: {
			Short: "ATPL",
			Long:  "Air Transport Pilot License",
		},
		FlightInstructorRating: {
			Short: "FI",
			Long:  "Flight Instructor",
		},
		FlightExaminerRating: {
			Short: "FE",
			Long:  "Flight Examiner",
		},
	}
)

func (r PilotRating) IsValidRating() bool {
	_, ok := pilotRatingMap[r]
	return ok
}

func (r PilotRating) Int() int {
	return int(r)
}

func (r PilotRating) Short() string {
	val, ok := pilotRatingMap[r]
	if ok {
		return val.Short
	}
	return "UNK"
}

func (r PilotRating) Long() string {
	val, ok := pilotRatingMap[r]
	if ok {
		return val.Long
	}
	return "Unknown"
}
