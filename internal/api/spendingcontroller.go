package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type SpendingController struct {
	TransactionRepository *models.TransactionRepository
}

//go:embed spendingbycategory.html.tmpl
var spendingByCategoryTmpl string

type spendingByCategory struct {
	Category        string
	Amount          string
	AmountHumanized string
}

type SpendingByCategoryDTO struct {
	AllTimeActive        bool
	Last12MonthsActive   bool
	Last6MonthsActive    bool
	Last3MonthsActive    bool
	SpendingByCategories []spendingByCategory
}

func (sc *SpendingController) spendingByCategoryHandler(w http.ResponseWriter, req *http.Request) {
	// Get time frame
	var relativeWindow int
	dto := SpendingByCategoryDTO{}
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

	// Get sum of all transactions in the time frame, grouped by category
	totalsByCategory, err := sc.TransactionRepository.GetSumOfTransactionsByCategory(startDate, endDate)
	if err != nil {
		fmt.Println("Error getting sum of transactions by category: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// Construct DTO
	for _, total := range totalsByCategory {
		dto.SpendingByCategories = append(dto.SpendingByCategories, spendingByCategory{
			Category:        total.Category,
			Amount:          utils.CentsToDollarStringMachineSafe(total.Amount),
			AmountHumanized: utils.CentsToDollarStringHumanized(total.Amount),
		})
	}

	// Render template
	tmpl, err := template.New("spendingByCategory").Parse(spendingByCategoryTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}
