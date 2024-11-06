package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/alexdglover/sage/internal/dependencyregistry"
	"gorm.io/gorm"
)

type Configuration struct {
	DbConnection *gorm.DB
}

var Config Configuration

func main() {

	ctx := context.TODO()
	dependencyRegistry := dependencyregistry.DependencyRegistry{}
	bootstrapper := dependencyRegistry.GetBootstrapper()
	bootstrapper.BootstrapDatabase(ctx)

	// open local browser to localhost:8080
	// openbrowser("http://localhost:8080")

	// start the API server
	apiServer, err := dependencyRegistry.GetApiServer()
	if err != nil {
		fmt.Println("Error starting API server")
		panic(err)
	}
	apiServer.StartApiServer(ctx)

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
