package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type NetIncomeController struct {
	TransactionRepository *models.TransactionRepository
}

type IncomeAndExpensesDataSet struct {
	Date      string
	Income    string
	Expenses  string
	NetIncome string
}

type IncomeAndExpensesDTO struct {
	DataSets []IncomeAndExpensesDataSet
}

//go:embed netincome.html.tmpl
var netIncomeTmpl string

func (nc *NetIncomeController) netIncomeHandler(w http.ResponseWriter, req *http.Request) {
	// We always start with today's date and work backwards based on relative window value
	endDate := time.Now()
	// Decrement relativeWindow by 1 (to account for the current month already being included)
	relativeWindow := 6
	relativeWindow = relativeWindow - 1
	// And calculate start date
	startDate := endDate.AddDate(0, (relativeWindow * -1), 0)

	// Get all transactions in asset accounts
	// Get all transactions from liability accounts and negate the amounts
	incomeTxns, err := nc.TransactionRepository.GetAllIncomeTransactions(context.TODO(), startDate, endDate)
	if err != nil {
		fmt.Println("error while getting asset transactions:", err)
		//TODO: add an HTTP return here
	}
	expenseTxns, err := nc.TransactionRepository.GetAllExpenseTransactions(context.TODO(), startDate, endDate)
	if err != nil {
		fmt.Println("error while getting asset transactions:", err)
		//TODO: add an HTTP return here
	}

	dto := IncomeAndExpensesDTO{}

	for idx, txnsWithDate := range incomeTxns {
		var runningIncomeTotal, runningExpenseTotal int64
		// Explicitly set these back to zero to avoid accumulating data across months
		runningIncomeTotal = 0
		runningExpenseTotal = 0

		for _, txn := range txnsWithDate.Transactions {
			runningIncomeTotal = runningIncomeTotal + txn.Amount
		}

		for _, txn := range expenseTxns[idx].Transactions {
			runningExpenseTotal = runningExpenseTotal - txn.Amount
		}
		netIncomeTotal := runningIncomeTotal + runningExpenseTotal

		incomeAndExpenses := IncomeAndExpensesDataSet{
			Date:      txnsWithDate.Date.Format("2006-01"),
			Income:    utils.CentsToDollarStringMachineSafe(runningIncomeTotal),
			Expenses:  utils.CentsToDollarStringMachineSafe(runningExpenseTotal),
			NetIncome: utils.CentsToDollarStringMachineSafe(netIncomeTotal),
		}
		dto.DataSets = append(dto.DataSets, incomeAndExpenses)
	}

	tmpl, err := template.New("netIncome").Parse(netIncomeTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}
