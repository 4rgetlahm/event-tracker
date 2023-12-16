package service

import (
	"errors"
	"fmt"

	"github.com/4rgetlahm/event-tracker/server/database"
	"github.com/4rgetlahm/event-tracker/server/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Slots       int    `json:"slots"`
}

func CreateEvent(req *EventCreateRequest) (entity.Event, error) {
	var db = database.GetDatabase()
	var event = entity.Event{
		Title:           req.Title,
		Description:     req.Description,
		Slots:           req.Slots,
		RegisteredUsers: []entity.User{},
	}
	err := db.Create(&event).Error
	return event, err
}

func UpdateEvent(UUID uuid.UUID, req *EventCreateRequest) (entity.Event, error) {
	var db = database.GetDatabase()
	var event entity.Event
	err := db.Find(&event, "id = ?", UUID).Error
	if err != nil {
		return event, err
	}
	event.Title = req.Title
	event.Description = req.Description
	event.Slots = req.Slots
	err = db.Save(&event).Error
	return event, err
}

func AddUserToEvent(UUID uuid.UUID, email string) (entity.Event, error) {
	var db = database.GetDatabase()

	var event entity.Event
	err := db.Model(entity.Event{}).Preload("RegisteredUsers").First(&event, "id = ?", UUID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, errors.New("Event not found")
		}
		return event, err
	}

	if event.Slots <= len(event.RegisteredUsers) {
		return event, errors.New("Event does not have enough slots")
	}

	if event.Status != string(entity.Open) {
		return event, errors.New("Event is not open")
	}

	fmt.Print(event)

	if emailInUsers(email, event.RegisteredUsers) {
		return event, errors.New("User is already registered to this event")
	}

	fmt.Print(event)

	db.Model(&event).Association("RegisteredUsers").Append(&entity.User{Email: email})

	return event, err
}

func RemoveUserFromEvent(UUID uuid.UUID, email string) (entity.Event, error) {
	var db = database.GetDatabase()

	var event entity.Event
	err := db.Model(entity.Event{}).Preload("RegisteredUsers").First(&event, "id = ?", UUID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, errors.New("Event not found")
		}
		return event, err
	}

	if event.Status != string(entity.Open) {
		return event, errors.New("Event is not open")
	}

	user, err := getUserByEmail(email, event.RegisteredUsers)
	if err != nil {
		return event, errors.New("User is not registered to this event")
	}

	db.Model(&event).Association("RegisteredUsers").Delete(&user)

	return event, err
}

func GetEvent(UUID uuid.UUID) (entity.Event, error) {
	var db = database.GetDatabase()
	var event entity.Event
	err := db.Model(entity.Event{}).Preload("RegisteredUsers").Find(&event, "id = ?", UUID).Error
	return event, err
}

func GetEvents(from int, to int) ([]entity.Event, error) {
	var db = database.GetDatabase()
	var events []entity.Event
	err := db.Limit(to - from).Offset(from).Find(&events).Error
	return events, err
}

func getUserByEmail(email string, users []entity.User) (entity.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}
	return entity.User{}, errors.New("User not found")
}

func emailInUsers(email string, users []entity.User) bool {
	fmt.Println(email)
	for _, user := range users {
		fmt.Println(user.Email)
		if user.Email == email {
			return true
		}
	}
	return false
}
