package api

import (
	_ "embed"
	"encoding/json"
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

	incomeJSON, _ := json.Marshal(cashFlowData.Income)
	expensesJSON, _ := json.Marshal(cashFlowData.Expenses)

	data := map[string]interface{}{
		"ActivePage":         "cashFlow",
		"Last3MonthsActive":  relativeWindow == "3",
		"Last6MonthsActive":  relativeWindow == "6",
		"Last12MonthsActive": relativeWindow == "12",
		"AllTimeActive":      relativeWindow == "allTime",
		"IncomeJSON":         string(incomeJSON),
		"ExpensesJSON":       string(expensesJSON),
	}

	tmpl := template.Must(template.New("cashFlow").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(cashFlowTmpl))
	err = utils.RenderTemplateAsHTML(w, tmpl, data)
	if err != nil {
		panic(err)
	}
}
