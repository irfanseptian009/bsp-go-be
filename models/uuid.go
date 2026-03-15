package models

import "github.com/google/uuid"

func newUUID() string {
	return uuid.NewString()
}
