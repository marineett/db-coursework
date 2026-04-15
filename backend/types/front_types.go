package types

type Category struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Subcategories []Subcategory `json:"subcategories"`
}

type Subcategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Feed struct {
	Repetitors []RepetitorData `json:"repetitors"`
	Categories []Category      `json:"categories"`
}

type RegistrationInfo struct {
	PersonalData
	AuthData
}
