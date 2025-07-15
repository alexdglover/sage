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
	ActivePage     string
	RelativeWindow string
	TotalIncome    string
	TotalExpenses  string
	Savings        string
	ShowSavings    bool
	Expenses       []ExpenseData
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
	savings := cashFlowData.TotalIncome - cashFlowData.TotalExpenses
	if savings < 0 {
		savings = 0 // Ensure savings is not negative
	}
	dto := CashFlowDto{
		ActivePage:     "cashflow",
		RelativeWindow: relativeWindow,
		TotalIncome:    utils.CentsToDollarStringMachineSafe(cashFlowData.TotalIncome),
		TotalExpenses:  utils.CentsToDollarStringMachineSafe(cashFlowData.TotalExpenses),
		Savings:        utils.CentsToDollarStringMachineSafe(savings),
		ShowSavings:    savings > 0,
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
