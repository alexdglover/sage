package api

import (
	_ "embed"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/services"
	"github.com/alexdglover/sage/internal/utils"
)

//go:embed cashflow.html
var cashFlowTmpl string

type CashFlowReportHandler struct {
	cashFlowService *services.CashFlowService
}

type ExpenseData struct {
	Name                string
	Amount              string
	AmountHumanReadable string
}

type CashFlowDto struct {
	ActivePage                 string
	NetIncome                  string
	NetIncomeHumanReadable     string
	NetIncomeLabel             string
	TotalIncome                string
	TotalIncomeHumanReadable   string
	TotalExpenses              string
	TotalExpensesHumanReadable string
	Expenses                   []ExpenseData
	StartDate                  string // YYYY-MM-DD for date input
	EndDate                    string // YYYY-MM-DD for date input
}

func NewCashFlowReportHandler(cfs *services.CashFlowService) *CashFlowReportHandler {
	return &CashFlowReportHandler{
		cashFlowService: cfs,
	}
}

func (h *CashFlowReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse startDate and endDate from query params
	query := r.URL.Query()
	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	var startDate, endDate time.Time
	var err error

	if endDateStr != "" {
		endDate = utils.ISO8601DateStringToTime(endDateStr)
	} else {
		endDate = time.Now()
		endDateStr = utils.TimeToISO8601DateString(endDate)
	}

	if startDateStr != "" {
		startDate = utils.ISO8601DateStringToTime(startDateStr)
	} else {
		startDate = endDate.AddDate(0, -3, 0)
		startDateStr = utils.TimeToISO8601DateString(startDate)
	}

	cashFlowData, err := h.cashFlowService.GetCashFlowData(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	netIncome := cashFlowData.TotalIncome - cashFlowData.TotalExpenses
	var netIncomeLabel string
	if netIncome < 0 {
		netIncome = -netIncome // Make it positive for display
		netIncomeLabel = "Taken From Savings"
	} else {
		netIncomeLabel = "Net Income"
	}

	dto := CashFlowDto{
		ActivePage:                 "cashflow",
		NetIncome:                  utils.CentsToDollarStringMachineSafe(netIncome),
		NetIncomeHumanReadable:     utils.CentsToDollarStringHumanized(netIncome),
		NetIncomeLabel:             netIncomeLabel,
		TotalIncome:                utils.CentsToDollarStringMachineSafe(cashFlowData.TotalIncome),
		TotalIncomeHumanReadable:   utils.CentsToDollarStringHumanized(cashFlowData.TotalIncome),
		TotalExpenses:              utils.CentsToDollarStringMachineSafe(cashFlowData.TotalExpenses),
		TotalExpensesHumanReadable: utils.CentsToDollarStringHumanized(cashFlowData.TotalExpenses),
		StartDate:                  startDateStr,
		EndDate:                    endDateStr,
	}
	for _, expense := range cashFlowData.Expenses {
		dto.Expenses = append(dto.Expenses, ExpenseData{
			Name:                expense.Name,
			Amount:              utils.CentsToDollarStringMachineSafe(expense.Amount),
			AmountHumanReadable: utils.CentsToDollarStringHumanized(expense.Amount),
		})
	}

	tmpl := template.Must(template.New("cashFlow").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(cashFlowTmpl))
	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}
