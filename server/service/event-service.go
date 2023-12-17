package service

import (
	"context"
	"errors"
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
		return event, err
	}

	return event, err
}

func GetEvents(from int, to int) ([]entity.Event, error) {
	var events []entity.Event
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cursor, err := database.GetEventCollection().Find(ctx, bson.M{}, options.Find().SetSkip(int64(from)).SetLimit(int64(to-from)))

	if err != nil {
		return events, err
	}

	if err = cursor.All(ctx, &events); err != nil {
		return events, err
	}

	return events, nil
}

func CreateEvent(req *EventCreateRequest) (entity.Event, error) {

	eventDate, err := time.Parse("2006-01-02", req.EventDate)

	if err != nil {
		return entity.Event{}, errors.New("Invalid event date format")
	}

	var event = entity.Event{
		ID:              primitive.NewObjectID(),
		Title:           req.Title,
		Description:     req.Description,
		Slots:           req.Slots,
		Status:          string(entity.Open),
		RegisteredUsers: []string{},
		Creator:         "test",
		EventDate:       eventDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = database.GetEventCollection().InsertOne(ctx, event)

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

	if emailInUsers(email, event.RegisteredUsers) {
		return event, errors.New("User is already registered to this event")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addQuery := bson.M{"$push": bson.M{"registeredusers": email}}
	filter := bson.M{"_id": objectId}

	_, err = database.GetEventCollection().UpdateOne(ctx, filter, addQuery)

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
