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
	Date                      string
	Income                    string
	IncomeHumanized           string
	Expenses                  string
	ExpensesHumanized         string
	NetIncome                 string
	NetIncomeHumanized        string
	TTMAverage                string
	TTMSeventyFifthPercentile string
	TTMTwentyFifthPercentile  string
}

type IncomeAndExpensesDTO struct {
	AllTimeActive      bool
	Last12MonthsActive bool
	Last6MonthsActive  bool
	Last3MonthsActive  bool
	DataSets           []IncomeAndExpensesDataSet
}

//go:embed netincome.html.tmpl
var netIncomeTmpl string

// netIncomeHandler is the HTTP handler for the net income page
// TODO: Implement logic for relative date controls
func (nc *NetIncomeController) netIncomeHandler(w http.ResponseWriter, req *http.Request) {

	dto := IncomeAndExpensesDTO{}
	var relativeWindow int

	if req.FormValue("relativeWindow") == "" {
		relativeWindow = 6
		dto.Last6MonthsActive = true
	} else {
		switch req.FormValue("relativeWindow") {
		case "12":
			relativeWindow = 12
			dto.Last12MonthsActive = true
		case "6":
			relativeWindow = 6
			dto.Last6MonthsActive = true
		case "3":
			relativeWindow = 3
			dto.Last3MonthsActive = true
		case "allTime":
			// 10 years in months, as a useful approximation of all time for now
			relativeWindow = 120
			dto.AllTimeActive = true
		default:
			fmt.Println("invalid relative window provided, falling back to 6 months")
			relativeWindow = 6
			dto.Last6MonthsActive = true
		}
	}

	// We always start with today's date and work backwards based on relative window value
	endDate := time.Now()
	relativeWindow = relativeWindow - 1
	// And calculate start date
	startDate := endDate.AddDate(0, (relativeWindow * -1), 0)

	// Get all income transactions
	incomeTxns, err := nc.TransactionRepository.GetAllIncomeTransactions(context.TODO(), startDate, endDate)
	if err != nil {
		fmt.Println("error while getting asset transactions:", err)
		//TODO: add an HTTP return here
	}
	// Get all expense transactions
	expenseTxns, err := nc.TransactionRepository.GetAllExpenseTransactions(context.TODO(), startDate, endDate)
	if err != nil {
		fmt.Println("error while getting asset transactions:", err)
		//TODO: add an HTTP return here
	}

	for idx, txnsWithDate := range incomeTxns {
		var runningIncomeTotal, runningExpenseTotal int
		runningIncomeTotal = 0
		runningExpenseTotal = 0

		for _, txn := range txnsWithDate.Transactions {
			runningIncomeTotal = runningIncomeTotal + txn.Amount
		}

		for _, txn := range expenseTxns[idx].Transactions {
			runningExpenseTotal = runningExpenseTotal - txn.Amount
		}
		netIncomeTotal := runningIncomeTotal + runningExpenseTotal

		netIncomeTTMAverage, twentyFifthPercentile, seventyFifthPercentile, _ := nc.TransactionRepository.GetTTMStatistics(context.TODO(), txnsWithDate.Date)

		incomeAndExpenses := IncomeAndExpensesDataSet{
			Date:                      txnsWithDate.Date.Format("2006-01"),
			Income:                    utils.CentsToDollarStringMachineSafe(runningIncomeTotal),
			IncomeHumanized:           utils.CentsToDollarStringHumanized(runningIncomeTotal),
			Expenses:                  utils.CentsToDollarStringMachineSafe(runningExpenseTotal),
			ExpensesHumanized:         utils.CentsToDollarStringHumanized(runningExpenseTotal),
			NetIncome:                 utils.CentsToDollarStringMachineSafe(netIncomeTotal),
			NetIncomeHumanized:        utils.CentsToDollarStringHumanized(netIncomeTotal),
			TTMAverage:                utils.CentsToDollarStringMachineSafe(netIncomeTTMAverage),
			TTMSeventyFifthPercentile: utils.CentsToDollarStringMachineSafe(seventyFifthPercentile),
			TTMTwentyFifthPercentile:  utils.CentsToDollarStringMachineSafe(twentyFifthPercentile),
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
