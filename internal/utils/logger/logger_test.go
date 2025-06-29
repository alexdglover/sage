package logger_test

import (
	"bufio"
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
	buf := &bytes.Buffer{}

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("Failed to redirect STDOUT")
	}

	os.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			buf.WriteString(scanner.Text())
		}
	}()

	logger := logger.Get()
	logger.Error("Error log")

	// Test output
	t.Log(buf)
	if buf.Len() == 0 {
		t.Error("No information logged to STDOUT")
	}
	if strings.Count(buf.String(), "\n") > 1 {
		t.Error("Expected only a single line of log output")
	}

	var logMessage logMesaage
	err = json.Unmarshal(buf.Bytes(), &logMessage)
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
		w.Close()
	})
}
