package types

type Department struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type DepartmentInitData struct {
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type HireInfo struct {
	ID           int64 `json:"id"`
	UserID       int64 `json:"user_id"`
	DepartmentID int64 `json:"department_id"`
}

type CompleteDepartmentInfo struct {
	Department
	Moderators []MoreratorProfileWithID
}
