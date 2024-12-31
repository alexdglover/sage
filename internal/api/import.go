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
	"github.com/alexdglover/sage/internal/utils"
)

type ImportController struct {
	AccountManager        *services.AccountManager
	ImportService         *services.ImportService
	TransactionRepository *models.TransactionRepository
}

//go:embed importStatementForm.html.tmpl
var importStatementFormTmpl string

//go:embed importStatusPage.html.tmpl
var importStatusPageTmpl string

type ImportStatementFormDTO struct {
	AccountNamesAndIDs []services.AccountNameAndID
}

type ImportStatusPageDTO struct {
	Submission   *models.ImportSubmission
	Transactions []TransactionDTO
}

func (ic *ImportController) importStatementFormHandler(w http.ResponseWriter, req *http.Request) {
	formDTO := ImportStatementFormDTO{}
	data, err := ic.AccountManager.GetAccountNamesAndIDs()
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

func (ic *ImportController) importSubmissionHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1024 * 1024 * 1024 * 4) // limit max input length to 4 GB

	file, header, err := req.FormFile("statementFile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := header.Filename
	accountID, err := utils.StringToUint(req.FormValue("accountSelector"))
	if err != nil {
		http.Error(w, "Unable to parse account ID", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	io.Copy(&buf, file)
	statement := buf.String()
	buf.Reset()

	// call service class to execute import
	importSubmission, err := ic.ImportService.ImportStatement(fileName, statement, accountID)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to import statement: %v", err)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	transactions, err := ic.TransactionRepository.GetTransactionsByImportSubmission(importSubmission.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get transactions for import submission: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	transactionDTOs := []TransactionDTO{}
	for _, txn := range transactions {
		transactionDTOs = append(transactionDTOs, TransactionDTO{
			ID:                 txn.ID,
			Date:               txn.Date,
			Description:        txn.Description,
			Amount:             utils.CentsToDollarStringHumanized(txn.Amount),
			Excluded:           txn.Excluded,
			AccountName:        txn.Account.Name,
			CategoryName:       txn.Category.Name,
			ImportSubmissionID: utils.UintPointerToString(txn.ImportSubmissionID),
		})
	}
	dto := ImportStatusPageDTO{
		Submission:   importSubmission,
		Transactions: transactionDTOs,
	}

	ic.importStatusHandler(w, dto)
}

// Handler to return HTML for the status of a single import submission
func (ic *ImportController) importStatusHandler(w http.ResponseWriter, dto ImportStatusPageDTO) {
	tmpl, err := template.New("importStatusPage").Parse(importStatusPageTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}
