package constants

type FacilityID string

const (
	// Special Facilities
	AcademyFacility      FacilityID = "ZAE"
	HeadquartersFacility FacilityID = "ZHQ"
	NonMemberFacility    FacilityID = "ZZN"
	InactiveFacility     FacilityID = "ZZI"

	// ARTCC Facilities
	AlbuquerqueFacility  FacilityID = "ZAB"
	AnchorageFacility    FacilityID = "ZAN"
	AtlantaFacility      FacilityID = "ZTL"
	BostonFacility       FacilityID = "ZBW"
	ChicagoFacility      FacilityID = "ZAU"
	ClevelandFacility    FacilityID = "ZOB"
	DenverFacility       FacilityID = "ZDV"
	FortWorthFacility    FacilityID = "ZFW"
	HonoluluFacility     FacilityID = "HCF"
	HoustonFacility      FacilityID = "ZHU"
	IndianapolisFacility FacilityID = "ZID"
	JacksonvilleFacility FacilityID = "ZJX"
	KansasCityFacility   FacilityID = "ZKC"
	LosAngelesFacility   FacilityID = "ZLA"
	MemphisFacility      FacilityID = "ZME"
	MiamiFacility        FacilityID = "ZMA"
	MinneapolisFacility  FacilityID = "ZMP"
	NewYorkFacility      FacilityID = "ZNY"
	OaklandFacility      FacilityID = "ZOA"
	SaltLakeFacility     FacilityID = "ZLC"
	SeattleFacility      FacilityID = "ZSE"
	WashingtonFacility   FacilityID = "ZDC"
)

var (
	FacilityDisplayNameMap = map[FacilityID]string{
		AcademyFacility:      "Academy",
		HeadquartersFacility: "Headquarters",
		NonMemberFacility:    "Non-Member",
		InactiveFacility:     "Inactive",

		AlbuquerqueFacility:  "Albuquerque ARTCC",
		AnchorageFacility:    "Anchorage ARTCC",
		AtlantaFacility:      "Atlanta ARTCC",
		BostonFacility:       "Boston ARTCC",
		ChicagoFacility:      "Chicago ARTCC",
		ClevelandFacility:    "Cleveland ARTCC",
		DenverFacility:       "Denver ARTCC",
		FortWorthFacility:    "Fort Worth ARTCC",
		HonoluluFacility:     "Honolulu Control FacilityID",
		HoustonFacility:      "Houston ARTCC",
		IndianapolisFacility: "Indianapolis ARTCC",
		JacksonvilleFacility: "Jacksonville ARTCC",
		KansasCityFacility:   "Kansas City ARTCC",
		LosAngelesFacility:   "Los Angeles ARTCC",
		MemphisFacility:      "Memphis ARTCC",
		MiamiFacility:        "Miami ARTCC",
		MinneapolisFacility:  "Minneapolis ARTCC",
		NewYorkFacility:      "New York ARTCC",
		OaklandFacility:      "Oakland ARTCC",
		SaltLakeFacility:     "Salt Lake ARTCC",
		SeattleFacility:      "Seattle ARTCC",
		WashingtonFacility:   "Washington, D.C. ARTCC",
	}
)

func (fac FacilityID) IsValidFacility() bool {
	_, ok := FacilityDisplayNameMap[fac]
	return ok
}

func (f FacilityID) DisplayName() string {
	return FacilityDisplayNameMap[f]
}
