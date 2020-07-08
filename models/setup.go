package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // required
)

// DB variable
var DB *gorm.DB

// ConnectDataBase function to establish database connection and handle migrations
func ConnectDataBase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{}, &Post{}, &Tag{})

	database.Create(&User{Username: "admin", Password: "123"})

	DB = database
}
