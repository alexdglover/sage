package api

import (
	"context"
	_ "embed"
	"log"
	"net/http"
)

//go:embed main.html.tmpl
var mainPageTmpl string

func StartApiServer(ctx context.Context) {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/net-worth", netWorthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
