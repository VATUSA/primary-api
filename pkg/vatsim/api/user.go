package vatsim_api

type User struct {
	CID      string `json:"cid"`
	Personal struct {
		FirstName string `json:"name_first"`
		LastName  string `json:"name_last"`
		FullName  string `json:"name_full"`
		Email     string `json:"email"`
	} `json:"personal"`
	VATSIM struct {
		ControllerRating struct {
			ID    int    `json:"id"`
			Short string `json:"short"`
			Long  string `json:"long"`
		} `json:"rating"`
		PilotRating struct {
			ID    int    `json:"id"`
			Short string `json:"short"`
			Long  string `json:"long"`
		} `json:"pilotrating"`
		Region struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"region"`
		Division struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"division"`
		Subdivision struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"subdivision"`
	}
}
