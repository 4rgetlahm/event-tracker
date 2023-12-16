package service

import (
	"github.com/4rgetlahm/event-tracker/server/entity"
	"github.com/google/uuid"
)

type EventCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Slots       int    `json:"slots"`
}

func CreateEvent(req *EventCreateRequest) (entity.Event, error) {
	return entity.Event{}, nil
}

func GetEvent(UUID uuid.UUID) (entity.Event, error) {
	return entity.Event{}, nil
}
