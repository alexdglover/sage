package main

import (
	"context"

	"gorm.io/gorm"

	"github.com/alexdglover/sage/internal/api"
	"github.com/alexdglover/sage/internal/models"
)

type Configuration struct {
	DbConnection *gorm.DB
}

var Config Configuration

func main() {

	ctx := context.TODO()

	models.BootstrapDatabase(ctx)

	// start an API server
	api.StartApiServer(ctx)

	// open local browser to localhost:8080

}
