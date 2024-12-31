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
	// And calculate start date
	startDate := endDate.AddDate(0, (relativeWindow * -1), 0)

	netIncomeDataByDate, err := nc.TransactionRepository.GetNetIncomeTotalsByDate(context.TODO(), startDate, endDate)
	if err != nil {
		// TODO: handle correctly
		fmt.Println("error while getting net income data:", err)
	}

	for _, netIncomeData := range netIncomeDataByDate {
		netIncomeTTMAverage, twentyFifthPercentile, seventyFifthPercentile, _ := nc.TransactionRepository.GetTTMStatistics(context.TODO(), netIncomeData.Date)

		incomeAndExpenses := IncomeAndExpensesDataSet{
			Date:                      netIncomeData.Date.Format("2006-01"),
			Income:                    utils.CentsToDollarStringMachineSafe(netIncomeData.Income),
			IncomeHumanized:           utils.CentsToDollarStringHumanized(netIncomeData.Income),
			Expenses:                  utils.CentsToDollarStringMachineSafe(netIncomeData.Expenses * -1),
			ExpensesHumanized:         utils.CentsToDollarStringHumanized(netIncomeData.Expenses * -1),
			NetIncome:                 utils.CentsToDollarStringMachineSafe(netIncomeData.NetIncome),
			NetIncomeHumanized:        utils.CentsToDollarStringHumanized(netIncomeData.NetIncome),
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
