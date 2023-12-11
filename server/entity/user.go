package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	uuid       uuid.UUID `json:"uuid"`
	email      string    `json:"email"`
	createDate time.Time `json:"createDate"`
}
