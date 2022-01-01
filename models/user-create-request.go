package models

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type UserCreateRequest struct {
	Email string `json:"email"`
}

func (ucr UserCreateRequest) Validate() error {
	return validation.ValidateStruct(
		&ucr,
		validation.Field(
			&ucr.Email,
			validation.Required,
			is.Email))
}
