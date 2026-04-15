package types

type InitModeratorData struct {
	InitUserData
}

type ModeratorData struct {
	UserData
	Salary      int64        `json:"salary"`
	Departments []Department `json:"departments"`
}

type ModeratorProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	MiddleName      string   `json:"middle_name"`
	TelephoneNumber string   `json:"telephone_number"`
	Email           string   `json:"email"`
	Salary          int64    `json:"salary"`
	Departments     []string `json:"departments"`
}

type MoreratorProfileWithID struct {
	ModeratorProfile
	ID int64 `json:"id"`
}
