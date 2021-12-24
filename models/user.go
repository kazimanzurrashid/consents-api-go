package models

type User struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Consents []Consent `json:"consents"`
}
