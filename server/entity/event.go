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
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Slots           int                `json:"slots"`
	RegisteredUsers []string           `json:"registeredUsers,omitempty"`
	Creator         string             `json:"creator"`
	Status          string             `json:"status"`
	EventDate       time.Time          `json:"eventDate"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

type EventWithUserState struct {
	Event
	IsRegistered bool `json:"isRegistered"`
	SlotsLeft    int  `json:"slotsLeft"`
}
