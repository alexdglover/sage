package main

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/alexdglover/sage/internal"
	"github.com/alexdglover/sage/internal/api"
)

func main() {

	ctx := context.TODO()

	db, err := gorm.Open(sqlite.Open("sage.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	internal.BootstrapTables(ctx, db)

	// start an API server
	api.StartApiServer(ctx)

	// open local browser to localhost:8080

}
