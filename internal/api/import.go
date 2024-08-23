package api

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

type ImportStatusPageDTO struct {
	Submission   *models.ImportSubmission
	Transactions []models.Transaction
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
	accountIDString := req.FormValue("accountSelector")
	accountID64, err := strconv.ParseUint(accountIDString, 10, 64)
	accountID := uint(accountID64)
	if err != nil {
		http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	io.Copy(&buf, file)
	statement := buf.String()
	buf.Reset()

	// call service class to execute import
	importSubmission, err := services.ImportStatement(fileName, statement, accountID)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to import statement: %v", err)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	tr := models.GetTransactionRepository()
	transactions, err := tr.GetTransactionsByImportSubmission(importSubmission.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get transactions for import submission: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	dto := ImportStatusPageDTO{
		Submission:   importSubmission,
		Transactions: transactions,
	}

	importStatusHandler(w, dto)
}

// Handler to return HTML for the status of a single import submission
func importStatusHandler(w http.ResponseWriter, dto ImportStatusPageDTO) {
	tmpl, err := template.New("importStatusPage").Parse(importStatusPageTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}
