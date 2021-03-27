package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	Connection *gorm.DB
)

func InitDatabase() {
	conn, err := gorm.Open(sqlite.Open("var/wg.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB: Could not open database: %s", err)
	}

	if err = conn.AutoMigrate(User{}, Client{}); err != nil {
		log.Fatalf("DB: Could not migrate database: %s", err)
	}

	Connection = conn
	log.Print("DB: configured.")
}
