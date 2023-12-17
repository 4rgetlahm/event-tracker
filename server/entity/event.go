package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventStatus string

const (
	Open   EventStatus = "open"
	Closed EventStatus = "closed"
)

type Event struct {
	ID              primitive.ObjectID `bson:"_id"`
	Title           string
	Description     string
	Slots           int
	RegisteredUsers []string
	Creator         string
	Status          string
	EventDate       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
