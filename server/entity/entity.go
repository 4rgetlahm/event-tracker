package entity

import "github.com/4rgetlahm/event-tracker/server/database"

func Init() {
	database.GetDatabase().AutoMigrate(&Event{})
}
