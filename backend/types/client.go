package types

type ClientData struct {
	UserData
	MeanRating float64  `json:"mean_rating"`
	Reviews    []Review `json:"reviews"`
}

type ClientProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	MiddleName      string   `json:"middle_name"`
	TelephoneNumber string   `json:"telephone_number"`
	Email           string   `json:"email"`
	MeanRating      float64  `json:"mean_rating"`
	Reviews         []Review `json:"reviews"`
}

type InitClientData struct {
	InitUserData
}
