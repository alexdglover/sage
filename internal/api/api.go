package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
)

func StartApiServer(ctx context.Context) {
	http.HandleFunc("/", dashboardHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	content, err := os.ReadFile("./internal/api/dashboard.html")
	if err != nil {
		log.Fatal(err)
	}
	contentString := string(content)
	io.WriteString(w, contentString)
}
