package api

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/utils"
)

type DashbaordController struct{}

//go:embed dashboard.html.tmpl
var dashboardTmpl string

// TODO: Consider moving this into a service class that returns just the data needed
func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	type emptyTmplVariables struct{}

	dto := emptyTmplVariables{}

	tmpl, err := template.New("dashboard").Parse(dashboardTmpl)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}
