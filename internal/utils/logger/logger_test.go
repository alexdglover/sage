package logger_test

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/alexdglover/sage/internal/utils/logger"
)

// This is a hack to move the working directory to the project root
// meaning the test created the log file in (Project Root)/logs
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

type logMesaage struct {
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
	Msg   string    `json:"msg"`
}

func TestErrorLog(t *testing.T) {
	// Save original stdout
	origStdout := os.Stdout

	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect stdout before logger is initialized
	os.Stdout = w

	logger := logger.Get()
	logger.Error("Error log")

	// Close writer to allow reading all output
	w.Close()
	os.Stdout = origStdout

	// Read output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()
	t.Log(output)
	if len(output) == 0 {
		t.Error("No information logged to STDOUT")
	}
	if strings.Count(output, "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}

	var logMessage logMesaage
	err = json.Unmarshal([]byte(output), &logMessage)
	if err != nil {
		t.Error("Error while decoding the log message", err.Error())
	}

	if logMessage.Level != "ERROR" {
		t.Errorf("Log level should be ERROR got: %s", logMessage.Level)
	}

	if logMessage.Msg != "Error log" {
		t.Errorf("Log message should be 'Error log' got: %s", logMessage.Msg)
	}

	//clean up log folder
	t.Cleanup(func() {
		cwd, _ := os.Getwd()
		if err := os.RemoveAll(path.Join(cwd, "/logs")); err != nil {
			log.Printf("ERROR: Removing test log folder: %v %s", err, cwd)
		}
	})
}
