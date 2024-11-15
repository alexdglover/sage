package api

import (
	_ "embed"
	"net/http"
	"text/template"
)

type DashbaordController struct{}

//go:embed dashboard.html.tmpl
var dashboardTmpl string

// TODO: Consider moving this into a service class that returns just the data needed
func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	type emptyTmplVariables struct{}

	foo := emptyTmplVariables{}

	tmpl, err := template.New("dashboard").Parse(dashboardTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, foo)
	if err != nil {
		panic(err)
	}
}
