package entity

import (
	"github.com/google/uuid"
)

type Event struct {
	uuid            uuid.UUID `json:"uuid"`
	title           string    `json:"title"`
	description     string    `json:"description"`
	slots           int       `json:"slots"`
	registeredUsers []User    `json:"registeredUsers"`
	creator         User      `json:"creator"`
}
