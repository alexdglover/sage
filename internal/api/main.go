package api

import (
	_ "embed"
	"net/http"
	"text/template"
)

type MainController struct{}

//go:embed main.html.tmpl
var mainTmpl string

// TODO: Consider moving this into a service class that returns just the data needed
func (mc *MainController) mainPageHandler(w http.ResponseWriter, req *http.Request) {
	type emptyTmplVariables struct{}

	foo := emptyTmplVariables{}

	tmpl, err := template.New("mainPage").Parse(mainPageTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, foo)
	if err != nil {
		panic(err)
	}
}
