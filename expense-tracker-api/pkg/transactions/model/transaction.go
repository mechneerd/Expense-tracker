package model

type Transaction struct {
	ID          string  `json:"id"`
	FamilyID    string  `json:"family_id"`
	UserID      *string `json:"user_id,omitempty"`
	Type        string  `json:"type"`
	Category    *string `json:"category,omitempty"`
	PaymentMethod *string `json:"payment_method,omitempty"`
	UPIApp      *string `json:"upi_app,omitempty"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}