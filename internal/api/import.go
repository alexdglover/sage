package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/services"
)

//go:embed importStatementForm.html.tmpl
var importStatementFormTmpl string

//go:embed importStatusPage.html.tmpl
var importStatusPageTmpl string

type ImportStatementFormDTO struct {
	AccountNamesAndIDs []services.AccountNameAndID
}

func importStatementFormHandler(w http.ResponseWriter, req *http.Request) {
	formDTO := ImportStatementFormDTO{}
	data, err := services.GetAccountNamesAndIDs()
	if err != nil {
		http.Error(w, "Unable to get account names and IDs", http.StatusInternalServerError)
	}
	formDTO.AccountNamesAndIDs = data

	// TODO: Get list of parsers to populate drop down

	tmpl, err := template.New("importStatementForm").Parse(importStatementFormTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, formDTO)
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
	accountId := req.FormValue("accountSelector")

	var buf bytes.Buffer
	io.Copy(&buf, file)
	statement := buf.String()
	buf.Reset()

	// call service class to execute import
	importSubmission, err := services.ImportStatement(fileName, statement, 1, "schwab")
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to import statement: %v", err)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	importStatusHandler(w, importSubmission)
}

// Handler to return HTML for the status of a single import submission
func importStatusHandler(w http.ResponseWriter, submission *models.ImportSubmission) {
	tmpl, err := template.New("importStatusPage").Parse(importStatusPageTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, *submission)
	if err != nil {
		panic(err)
	}
}
