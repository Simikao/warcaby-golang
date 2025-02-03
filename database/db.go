package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nick     string `gorm:"uniqueIndex; not null"`
	Email    string `gorm:"uniqueIndex; not null"`
	Bio      string
	Password string `gorm:"not null"`
}

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
