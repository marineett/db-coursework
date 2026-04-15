package types

import "time"

type ContractStatus int

const (
	ContractStatusNull ContractStatus = iota
	ContractStatusPending
	ContractStatusActive
	ContractStatusCompleted
	ContractStatusCancelled
	ContractStatusBanned
)

type PaymentStatus int

const (
	PaymentStatusNull PaymentStatus = iota
	PaymentStatusPending
	PaymentStatusPaid
	PaymentStatusRefused
	PaymentStatusRefunded
)

type ContractCategory int

const (
	ContractCategoryNull ContractCategory = iota
	ContractCategoryTranslation
	ContractCategoryWriting
	ContractCategoryDesign
	ContractCategoryProgramming
	ContractCategoryOther
)

type ContractSubcategory int

const (
	ContractSubcategoryNull ContractSubcategory = iota
	ContractSubcategoryTutoring
	ContractSubcategoryTranslation
	ContractSubcategoryWriting
	ContractSubcategoryDesign
	ContractSubcategoryProgramming
	ContractSubcategoryOther
)

type ContractInitInfo struct {
	ClientID              int64                 `json:"client_id"`
	ContractCategory      ContractCategory      `json:"contract_category"`
	ContractSubcategories []ContractSubcategory `json:"contract_subcategories"`
	Description           string                `json:"description"`
	Price                 int64                 `json:"price"`
	Commission            int64                 `json:"commission"`
	StartDate             time.Time             `json:"start_date"`
	Duration              int64                 `json:"duration"`
}

type ContractUpdateInfo struct {
	Status            ContractStatus `json:"status"`
	PriceChange       int64          `json:"price_change"`
	DescriptionChange string         `json:"description_change"`
}

type Contract struct {
	ID                int64          `json:"id"`
	ClientID          int64          `json:"client_id"`
	RepetitorID       int64          `json:"repetitor_id"`
	TransactionID     int64          `json:"transaction_id"`
	CreatedAt         time.Time      `json:"created_at"`
	Description       string         `json:"description"`
	Status            ContractStatus `json:"status"`
	PaymentStatus     PaymentStatus  `json:"payment_status"`
	ReviewClientID    int64          `json:"review_client_id"`
	ReviewRepetitorID int64          `json:"review_repetitor_id"`
	Price             int64          `json:"price"`
	Commission        int64          `json:"commission"`
	StartDate         time.Time      `json:"start_date"`
	EndDate           time.Time      `json:"end_date"`
	IDCRChat          int64          `json:"id_client_repetitor_chat"`
	IDCMRepChat       int64          `json:"id_client_manager_chat"`
	IDRMRepChat       int64          `json:"id_repetitor_moderator_chat"`
}

type Review struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}
