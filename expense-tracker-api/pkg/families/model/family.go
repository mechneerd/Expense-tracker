package model

type Family struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UniqueCode string `json:"unique_code"`
	CreatedBy string `json:"created_by"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type FamilyMember struct {
	ID        string `json:"id"`
	FamilyID string `json:"family_id"`
	UserID    string `json:"user_id"`
	FamilyRole string `json:"family_role"`
	Status    string `json:"status"`
	JoinedAt  string `json:"joined_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}