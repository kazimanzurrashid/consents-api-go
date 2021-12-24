package services

import (
	"strings"

	"github.com/google/uuid"
)

func generateID() string {
	return strings.ToLower(uuid.New().String())
}
