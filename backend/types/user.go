package types

import "time"

type UserType int

const (
	Unauthorized UserType = iota
	Client
	Repetitor
	Moderator
	Admin
)

type InitUserData struct {
	PersonalData PersonalData `json:"personal_data"`
	AuthData     AuthData     `json:"auth_data"`
}

type UserData struct {
	ID               int64     `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLoginDate    time.Time `json:"last_login_date"`
	PersonalDataID   int64     `json:"personal_data_id"`
}

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthVerdict struct {
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
}

type AuthInfo struct {
	ID       int64    `json:"id"`
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
	Login    string   `json:"login"`
	Password string   `json:"password"`
}

type PassportData struct {
	PassportNumber   string    `json:"passport_number"`
	PassportSeries   string    `json:"passport_series"`
	PassportDate     time.Time `json:"passport_date"`
	PassportIssuedBy string    `json:"passport_issued_by"`
}

type PersonalData struct {
	ID              int64  `json:"id"`
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	PassportData
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type RegistrationData struct {
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	UserType         UserType  `json:"user_type"`
	RegistrationDate time.Time `json:"registration_date"`
}
