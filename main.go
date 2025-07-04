package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/alexdglover/sage/internal/dependencyregistry"
	"github.com/alexdglover/sage/internal/utils/logger"
)

func main() {

	logger := logger.Get()

	ctx := context.TODO()
	dependencyRegistry := dependencyregistry.DependencyRegistry{}
	bootstrapper := dependencyRegistry.GetBootstrapper()
	bootstrapper.BootstrapDatabase(ctx)

	// open local browser to localhost:8080 if the config is set to true
	settingsRepository, err := dependencyRegistry.GetSettingsRepository()
	if err != nil {
		logger.Error("Error while getting configRepository")
		panic(err)
	}
	settings, err := settingsRepository.GetSettings()
	if err != nil {
		logger.Error("Error while getting settings")
		panic(err)
	}
	if (*settings).LaunchBrowserOnStartup {
		openbrowser("http://localhost:8080")
	}

	// start the API server
	apiServer, err := dependencyRegistry.GetApiServer()
	if err != nil {

		logger.Error("Error starting API server")
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
