package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	req.ParseMultipartForm(32 << 20) // limit your max input length!

	// just for debugging, seeing form fields
	fmt.Println(req.MultipartForm)
	file, header, err := req.FormFile("statementFile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", fileName[0]+fileName[1])

	var buf bytes.Buffer
	io.Copy(&buf, file)
	contents := buf.String()
	fmt.Println(contents)
	buf.Reset()

	// Use the body (in this case, we're just printing it)
	fmt.Fprintf(w, "Received: %s", contents)
}
