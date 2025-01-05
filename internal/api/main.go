package api

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/utils"
)

type MainController struct{}

//go:embed main.html.tmpl
var mainPageTmpl string

// TODO: Consider moving this into a service class that returns just the data needed
func (mc *MainController) mainPageHandler(w http.ResponseWriter, req *http.Request) {
	type emptyTmplVariables struct{}

	dto := emptyTmplVariables{}

	tmpl, err := template.New("mainPage").Parse(mainPageTmpl)
	if err != nil {
		panic(err)
	}
	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}
