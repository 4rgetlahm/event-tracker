package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDatabase() *gorm.DB {
	if DB == nil {
		err := Connect()
		if err != nil {
			panic(err)
		}
	}
	return DB
}

func Connect() error {
	dsn := "host=localhost user=? password=? dbname=? port=? TimeZone=Europe/Vilnius"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = db
	return err
}
