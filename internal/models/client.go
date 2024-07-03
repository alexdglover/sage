package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Function to set the DB client to be used by all models and repositories
// TODO: Consider re-visiting this as an injected config instead of a package-wide
// variable
func createDbClient() {
	dbClient, err := gorm.Open(sqlite.Open("sage.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = dbClient
}
