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
	Name   string
	Amount string
}

type CashFlowDto struct {
	ActivePage                 string
	NetIncome                  string
	NetIncomeHumanReadable     string
	NetIncomeLabel             string
	RelativeWindow             string
	TotalIncome                string
	TotalIncomeHumanReadable   string
	TotalExpenses              string
	TotalExpensesHumanReadable string
	Savings                    string
	ShowSavings                bool
	Expenses                   []ExpenseData
}

func NewCashFlowReportHandler(cfs *services.CashFlowService) *CashFlowReportHandler {
	return &CashFlowReportHandler{
		cashFlowService: cfs,
	}
}

func (h *CashFlowReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	relativeWindow := r.URL.Query().Get("relativeWindow")

	var startDate, endDate time.Time
	endDate = time.Now()

	switch relativeWindow {
	case "3":
		startDate = endDate.AddDate(0, -3, 0)
	case "6":
		startDate = endDate.AddDate(0, -6, 0)
	case "12":
		startDate = endDate.AddDate(0, -12, 0)
	default:
		// Default to all time
		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		relativeWindow = "allTime"
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
		netIncomeLabel = "Net Loss"
	} else {
		netIncomeLabel = "Net Income"
	}

	dto := CashFlowDto{
		ActivePage:                 "cashflow",
		RelativeWindow:             relativeWindow,
		NetIncome:                  utils.CentsToDollarStringMachineSafe(netIncome),
		NetIncomeHumanReadable:     utils.CentsToDollarStringHumanized(netIncome),
		NetIncomeLabel:             netIncomeLabel,
		TotalIncome:                utils.CentsToDollarStringMachineSafe(cashFlowData.TotalIncome),
		TotalIncomeHumanReadable:   utils.CentsToDollarStringHumanized(cashFlowData.TotalIncome),
		TotalExpenses:              utils.CentsToDollarStringMachineSafe(cashFlowData.TotalExpenses),
		TotalExpensesHumanReadable: utils.CentsToDollarStringHumanized(cashFlowData.TotalExpenses),
	}
	for _, expense := range cashFlowData.Expenses {
		dto.Expenses = append(dto.Expenses, ExpenseData{
			Name:   expense.Name,
			Amount: utils.CentsToDollarStringMachineSafe(expense.Amount),
		})
	}

	tmpl := template.Must(template.New("cashFlow").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(cashFlowTmpl))
	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}
