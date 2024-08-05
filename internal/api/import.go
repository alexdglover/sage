package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

//go:embed importPage.html.tmpl
var importPageTmpl string

func importPageHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: Get list of accounts to populate account choice drop down

	// TODO: Get list of parsers to populate drop down

	tmpl, err := template.New("importStatementPage").Parse(importPageTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func importSubmissionHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1024 * 1024 * 1024 * 4) // limit max input length to 4 GB

	file, header, err := req.FormFile("statementFile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := header.Filename
	fmt.Println("fileName is", fileName)

	var buf bytes.Buffer
	io.Copy(&buf, file)
	contents := buf.String()
	buf.Reset()

	// call service class to execute import
	// result, err := importStatement(statement, accountID, parserID)
	// if err != nil {
	//	errorMessage := fmt.Sprintf("Unable to import statement: %v", err)
	// 	http.Error(w, errorMessage, http.StatusBadRequest)
	// }

	// Use the body (in this case, we're just printing it)
	fmt.Fprintf(w, "Received: %s", contents)
}
