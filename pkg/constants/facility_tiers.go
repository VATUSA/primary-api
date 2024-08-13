package constants

var (
	FacilityTierOne = map[FacilityID][]FacilityID{
		AlbuquerqueFacility: {
			LosAngelesFacility,
			DenverFacility,
			KansasCityFacility,
			FortWorthFacility,
			HoustonFacility,
		},
		LosAngelesFacility: {
			OaklandFacility,
			SaltLakeFacility,
			DenverFacility,
			AlbuquerqueFacility,
			HonoluluFacility,
		},
		OaklandFacility: {
			SeattleFacility,
			SaltLakeFacility,
			LosAngelesFacility,
			HonoluluFacility,
		},
		SeattleFacility: {
			AnchorageFacility,
			SaltLakeFacility,
			OaklandFacility,
			HonoluluFacility,
		},
		HonoluluFacility: {
			AnchorageFacility,
			SeattleFacility,
			OaklandFacility,
			LosAngelesFacility,
		},
		DenverFacility: {
			SaltLakeFacility,
			MinneapolisFacility,
			KansasCityFacility,
			AlbuquerqueFacility,
			LosAngelesFacility,
		},
		HoustonFacility: {
			AlbuquerqueFacility,
			FortWorthFacility,
			MemphisFacility,
			JacksonvilleFacility,
		},

		//TODO - finish me
	}
)

func (a FacilityID) IsTierOne(b FacilityID) bool {
	for _, fac := range FacilityTierOne[a] {
		if fac == b {
			return true
		}
	}

	return false
}
