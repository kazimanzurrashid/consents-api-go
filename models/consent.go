package models

import "github.com/go-ozzo/ozzo-validation"

type Consent struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

func (c Consent) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(
			&c.ID,
			validation.Required,
			validation.In(ConsentEmail, ConsentSMS)))
}
