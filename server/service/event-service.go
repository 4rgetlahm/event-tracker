package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/4rgetlahm/event-tracker/server/database"
	"github.com/4rgetlahm/event-tracker/server/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type EventCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Slots       int    `json:"slots"`
	EventDate   string `json:"eventDate"`
}

func GetEvent(objectId primitive.ObjectID) (entity.Event, error) {
	var event entity.Event
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := database.GetEventCollection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&event)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return event, errors.New("Event not found")
		}

		log.Println(err)
		return event, errors.New("Error retrieving event")
	}

	return event, err
}

func GetEventsForUser(email string, from int, to int) ([]entity.EventWithUserState, error) {
	if from < 0 {
		return []entity.EventWithUserState{}, errors.New("Invalid from field")
	}

	if to < 0 {
		return []entity.EventWithUserState{}, errors.New("Invalid to field")
	}

	if from > to {
		return []entity.EventWithUserState{}, errors.New("From cannot be greater than to")
	}

	if to-from > 100 {
		return []entity.EventWithUserState{}, errors.New("Cannot retrieve more than 100 events")
	}

	var events []entity.Event
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cursor, err := database.GetEventCollection().Find(ctx, bson.M{}, options.Find().
		SetSort(bson.D{{Key: "eventdate", Value: -1}}).
		SetSkip(int64(from)).
		SetLimit(int64(to-from)))

	if err != nil {
		log.Println(err)
		return []entity.EventWithUserState{}, errors.New("Error retrieving events")
	}

	if err = cursor.All(ctx, &events); err != nil {
		log.Println(err)
		return []entity.EventWithUserState{}, errors.New("Error retrieving events")
	}

	var detailedEvents []entity.EventWithUserState

	for _, event := range events {
		isRegistered := emailInUsers(email, event.RegisteredUsers)
		slotsLeft := event.Slots - len(event.RegisteredUsers)
		event.RegisteredUsers = nil
		event.Creator = "hidden"
		detailedEvents = append(detailedEvents, entity.EventWithUserState{
			Event:        event,
			IsRegistered: isRegistered,
			SlotsLeft:    slotsLeft,
		})
	}

	return detailedEvents, nil

}

func GetDetailedEvents(email string, from int, to int) ([]entity.EventWithUserState, error) {
	if from < 0 {
		return []entity.EventWithUserState{}, errors.New("Invalid from field")
	}

	if to < 0 {
		return []entity.EventWithUserState{}, errors.New("Invalid to field")
	}

	if from > to {
		return []entity.EventWithUserState{}, errors.New("From cannot be greater than to")
	}

	if to-from > 100 {
		return []entity.EventWithUserState{}, errors.New("Cannot retrieve more than 100 events")
	}

	var events []entity.Event
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cursor, err := database.GetEventCollection().Find(ctx, bson.M{}, options.Find().
		SetSort(bson.D{{Key: "eventdate", Value: -1}}).
		SetSkip(int64(from)).
		SetLimit(int64(to-from)))

	if err != nil {
		log.Println(err)
		return []entity.EventWithUserState{}, errors.New("Error retrieving events")
	}

	if err = cursor.All(ctx, &events); err != nil {
		log.Println(err)
		return []entity.EventWithUserState{}, errors.New("Error retrieving events")
	}

	var detailedEvents []entity.EventWithUserState

	for _, event := range events {
		detailedEvents = append(detailedEvents, entity.EventWithUserState{
			Event:        event,
			IsRegistered: emailInUsers(email, event.RegisteredUsers),
			SlotsLeft:    event.Slots - len(event.RegisteredUsers),
		})
	}

	return detailedEvents, nil
}

func CreateEvent(creatorEmail string, req *EventCreateRequest) (entity.Event, error) {

	eventDate, err := time.Parse("2006-01-02", req.EventDate)

	if err != nil {
		return entity.Event{}, errors.New("Invalid event date format")
	}

	if eventDate.Before(time.Now()) {
		return entity.Event{}, errors.New("Event date cannot be in the past")
	}

	var event = entity.Event{
		ID:              primitive.NewObjectID(),
		Title:           req.Title,
		Description:     req.Description,
		Slots:           req.Slots,
		Status:          string(entity.Open),
		RegisteredUsers: []string{},
		Creator:         creatorEmail,
		EventDate:       eventDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = database.GetEventCollection().InsertOne(ctx, event)

	if err != nil {
		log.Println(err)
		return entity.Event{}, errors.New("Error creating event")
	}

	return event, err
}

func UpdateEvent(objectId primitive.ObjectID, req *EventCreateRequest) (entity.Event, error) {
	event, err := GetEvent(objectId)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return event, errors.New("Event not found")
		}
		return event, err
	}

	event.Title = req.Title
	event.Description = req.Description
	event.Slots = req.Slots

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectId}
	_, err = database.GetEventCollection().UpdateOne(ctx, filter, event)

	if err != nil {
		log.Println(err)
		return event, errors.New("Error updating event")
	}

	return event, err
}

func AddUserToEvent(objectId primitive.ObjectID, email string) (entity.Event, error) {
	event, err := GetEvent(objectId)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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

	if event.EventDate.Before(time.Now()) {
		return event, errors.New("Event has already passed")
	}

	if emailInUsers(email, event.RegisteredUsers) {
		return event, errors.New("User is already registered to this event")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addQuery := bson.M{"$push": bson.M{"registeredusers": email}}
	filter := bson.M{"_id": objectId}

	_, err = database.GetEventCollection().UpdateOne(ctx, filter, addQuery)

	if err != nil {
		log.Println(err)
		return event, errors.New("Error adding user to event")
	}

	return event, err
}

func RemoveUserFromEvent(objectId primitive.ObjectID, email string) (entity.Event, error) {

	event, err := GetEvent(objectId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return event, errors.New("Event not found")
		}
		return event, err
	}

	if event.Status != string(entity.Open) {
		return event, errors.New("Event is not open")
	}

	if !emailInUsers(email, event.RegisteredUsers) {
		return event, errors.New("User is not registered to this event")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	removeQuery := bson.M{"$pull": bson.M{"registeredusers": email}}
	filter := bson.M{"_id": objectId}

	_, err = database.GetEventCollection().UpdateOne(ctx, filter, removeQuery)

	if err != nil {
		log.Println(err)
		return event, errors.New("Error removing user from event")
	}

	return event, err
}

func emailInUsers(email string, userEmails []string) bool {
	for _, userEmail := range userEmails {
		if userEmail == email {
			return true
		}
	}
	return false
}
