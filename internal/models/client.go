package models

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Function to set the DB client to be used by all models and repositories
// TODO: Consider re-visiting this as an injected config instead of a package-wide
// variable
func createDbClient() {

	// TODO: Reduce logger verbosity once stable
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Include params in the SQL log
			Colorful:                  true,        // Enable color
		},
	)

	var sageFilePath string
	if sageFileEnvVar, ok := os.LookupEnv("SAGE_FILE"); ok {
		sageFilePath = sageFileEnvVar
	} else {
		sageFilePath = "sage.db"
	}
	dbClient, err := gorm.Open(sqlite.Open(sageFilePath), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db = dbClient
}
