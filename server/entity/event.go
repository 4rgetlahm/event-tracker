package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatus string

const (
	Open   EventStatus = "open"
	Closed EventStatus = "closed"
)

type Event struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title           string
	Description     string
	Slots           int
	RegisteredUsers []User `gorm:"many2many:event_registrations;"`
	Creator         string `gorm:"not null"`
	Status          string `gorm:"default:'open'"`
	EventDate       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
