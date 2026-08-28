package model

import "time"

type User struct {
	ID              string    `json:"id"`
	GoogleID        *string   `json:"google_id,omitempty"`
	Email           string    `json:"email"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone,omitempty"`
	AvatarURL       *string   `json:"avatar_url,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewUser(email, firstName, lastName string) *User {
	return &User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}