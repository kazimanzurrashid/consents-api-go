package models

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EventCreateUser struct {
	ID string `json:"id"`
}

func (ecu EventCreateUser) Validate() error {
	return validation.ValidateStruct(
		&ecu,
		validation.Field(
			&ecu.ID,
			validation.Required,
			is.UUID))
}

type EventCreateRequest struct {
	User     *EventCreateUser `json:"user"`
	Consents *[]Consent       `json:"consents"`
}

func (ecr EventCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&ecr,
		validation.Field(&ecr.User, validation.Required),
		validation.Field(&ecr.Consents, validation.Required))
}
