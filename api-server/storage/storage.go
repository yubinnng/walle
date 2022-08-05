package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
var POSTGRES_USER = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")

var Client *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=walle port=5432", POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD)
	log.Println("Connecting PostgreSQL: " + dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database!")
	}
	Client = db
}
