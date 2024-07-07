package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

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

	// open local browser to localhost:8080
	openbrowser("http://localhost:8080")

	// start the API server
	api.StartApiServer(ctx)

}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
