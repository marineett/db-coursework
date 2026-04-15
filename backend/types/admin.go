package types

type InitAdminData struct {
	InitUserData
	Salary int64 `json:"salary"`
}

type AdminData struct {
	UserData
	DepartmentID int64 `json:"department_id"`
	Salary       int64 `json:"salary"`
}

type AdminProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	Salary          int64  `json:"salary"`
}
