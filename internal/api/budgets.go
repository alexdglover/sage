package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/services"
	"github.com/alexdglover/sage/internal/utils"
)

type BudgetController struct {
	BudgetRepository      *models.BudgetRepository
	BudgetService         *services.BudgetService
	CategoryRepository    *models.CategoryRepository
	TransactionRepository *models.TransactionRepository
}

//go:embed budgets.html
var budgetsPageTmpl string

//go:embed budgetForm.html
var budgetsFormTmpl string

//go:embed budgetDetail.html
var budgetDetailTmpl string

const ColorGreen string = "#198754"
const ColorRed string = "#dc3545"

type BudgetDTO struct {
	ID           uint
	CategoryName string
	Amount       string
	Spend        string
	PercentUsed  int
}

type BudgetsPageDTO struct {
	Budgets     []BudgetDTO
	BudgetSaved bool
}

type BudgetFormDTO struct {
	// If we're updating an existing budget in the form, Updating will be true
	// If we're creating a new budget, Updating will be false
	Updating     bool
	BudgetID     string
	CategoryName string
	Amount       string
	Categories   []models.Category
}

type BudgetDataByMonthDTO struct {
	Amount string
	Color  string
	Month  string
	Spend  string
}

type BudgetDetailDTO struct {
	ID                  uint
	CategoryName        string
	NumOfMonthsExceeded string
	ExceededColor       string
	Average             string
	StdDev              string
	Volatility          string
	BudgetData          []BudgetDataByMonthDTO
}

func (bc *BudgetController) generateBudgetForm(w http.ResponseWriter, req *http.Request) {
	var dto BudgetFormDTO

	categories, err := bc.CategoryRepository.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
	}
	dto.Categories = categories

	budgetIDQueryParameter := req.URL.Query().Get("budgetID")
	if budgetIDQueryParameter != "" {
		budgetID, err := utils.StringToUint(budgetIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusInternalServerError)
			return
		}
		budget, err := bc.BudgetRepository.GetBudgetByID(budgetID)
		if err != nil {
			http.Error(w, "Unable to get budget", http.StatusInternalServerError)
			return
		}

		dto.Updating = true
		dto.BudgetID = fmt.Sprint(budget.ID)
		dto.CategoryName = budget.Category.Name
		dto.Amount = utils.CentsToDollarStringHumanized(budget.Amount)
	} else {
		dto.Updating = false
	}

	// If a user creates a budget from the categories page, we should pre-populate the category dropdown even
	// though a budget doesn't exist yet
	categoryNameQueryParameter := req.URL.Query().Get("categoryName")
	if categoryNameQueryParameter != "" {
		dto.CategoryName = categoryNameQueryParameter
	}

	tmpl, err := template.New("budgetForm").Parse(budgetsFormTmpl)
	if err != nil {
		panic(err)
	}

	err = utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}

func (bc *BudgetController) upsertBudget(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	budgetID := req.FormValue("budgetID")
	budgetCategory := req.FormValue("budgetCategory")
	amount := req.FormValue("amount")

	var budget models.Budget

	if budgetID != "" {
		id, err := utils.StringToUint(budgetID)
		if err != nil {
			http.Error(w, "Unable to parse budget ID", http.StatusBadRequest)
			return
		}
		budget, err = bc.BudgetRepository.GetBudgetByID(id)
		if err != nil {
			http.Error(w, "Unable to get budget", http.StatusBadRequest)
			return
		}
	} else {
		budget = models.Budget{}
	}

	budgetCategoryID, err := utils.StringToUint(budgetCategory)
	if err != nil {
		http.Error(w, "Unable to parse category ID", http.StatusBadRequest)
		return
	}
	budget.CategoryID = budgetCategoryID
	category, err := bc.CategoryRepository.GetCategoryByID(budget.CategoryID)
	if err != nil {
		panic("failed to get category from CategoryID")
	}
	budget.Category = category

	budget.Amount = utils.DollarStringToCents(amount)

	_, err = bc.BudgetRepository.Save(budget)
	if err != nil {
		http.Error(w, "Unable to save budget", http.StatusBadRequest)
		return
	}

	bc.sendViewResponse(w, true)
}

func (bc *BudgetController) deleteBudget(w http.ResponseWriter, req *http.Request) {
	budgetIDInput := req.FormValue("budgetID")

	budgetID, err := utils.StringToUint(budgetIDInput)
	if err != nil {
		http.Error(w, "Unable to parse a budget ID from input", http.StatusBadRequest)
		return
	}

	err = bc.BudgetRepository.DeleteByID(budgetID)
	if err != nil {
		http.Error(w, "Unable to delete budget", http.StatusBadRequest)
		return
	}
	bc.sendViewResponse(w, true)
}

func (bc *BudgetController) generateBudgetsView(w http.ResponseWriter, req *http.Request) {
	bc.sendViewResponse(w, false)
}

// Generic function to send the view response
func (bc *BudgetController) sendViewResponse(w http.ResponseWriter, update bool) {
	// Get all budget and associated spend
	budgetsAndSpend, err := bc.BudgetService.GetAllBudgetsAndCurrentSpend()
	if err != nil {
		http.Error(w, "Unable to parse a budgets and spend", http.StatusBadRequest)
		return
	}

	// Build budgets DTO
	budgetsDTO := make([]BudgetDTO, len(budgetsAndSpend))
	for i, budget := range budgetsAndSpend {
		budgetsDTO[i] = BudgetDTO{
			ID:           budget.ID,
			CategoryName: budget.CategoryName,
			Amount:       utils.CentsToDollarStringHumanized(budget.Amount),
			Spend:        utils.CentsToDollarStringHumanized(budget.Spend),
			PercentUsed:  budget.PercentUsed,
		}
	}
	budgetsPageDTO := BudgetsPageDTO{
		Budgets: budgetsDTO,
	}
	if update {
		budgetsPageDTO.BudgetSaved = true
	}

	tmpl := template.Must(template.New("budgetsPage").Parse(pageComponents))
	tmpl = template.Must(tmpl.Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(budgetsPageTmpl))

	err = utils.RenderTemplateAsHTML(w, tmpl, budgetsPageDTO)
	if err != nil {
		panic(err)
	}
}

// Generate the budget details view, pulling detailed statistics about the budget over the last 6 months
func (bc *BudgetController) generateBudgetDetailsView(w http.ResponseWriter, req *http.Request) {
	budgetIDQueryParameter := req.URL.Query().Get("budgetID")
	budgetID, err := utils.StringToUint(budgetIDQueryParameter)
	if err != nil {
		message := fmt.Sprint("Unable to parse budget ID", budgetIDQueryParameter)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	budget, err := bc.BudgetRepository.GetBudgetByID(budgetID)
	if err != nil {
		message := fmt.Sprint("Unable to get budget", budgetID)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Get comparison of budget vs spend for the last 6 months
	spendAndBudgetByMonth, err := bc.BudgetService.GetBudgetAndMonthlySpend(budgetID, 6)
	if err != nil {
		message := fmt.Sprintf("Error fetching spend and budget info for budget: %v Reason: %v", budgetID, err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	timesExceeded := 0
	budgetData := []BudgetDataByMonthDTO{}
	for _, spendAndBudget := range spendAndBudgetByMonth {
		if spendAndBudget.Spend > spendAndBudget.Amount {
			timesExceeded++
		}
		var color string
		if spendAndBudget.Spend > spendAndBudget.Amount {
			color = ColorRed
		} else {
			color = ColorGreen
		}
		budgetData = append(budgetData, BudgetDataByMonthDTO{
			Amount: utils.CentsToDollarStringMachineSafe(spendAndBudget.Amount),
			Color:  color,
			Month:  utils.ConvertTimeToMonString(spendAndBudget.Month),
			Spend:  utils.CentsToDollarStringMachineSafe(spendAndBudget.Spend),
		})
	}

	averageSpend, spendStdDeviation, err := bc.BudgetService.GetMeanAndStandardDeviation(budgetID, 6)
	if err != nil {
		message := fmt.Sprintf("Error fetching average spend and standard deviation for budget: %v Reason: %v", budgetID, err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	budgetDetailDTO := BudgetDetailDTO{
		ID:                  budgetID,
		CategoryName:        budget.Category.Name,
		NumOfMonthsExceeded: fmt.Sprint(timesExceeded),
		Average:             utils.CentsToDollarStringHumanized(averageSpend),
		StdDev:              utils.CentsToDollarStringHumanized(spendStdDeviation),
		BudgetData:          budgetData,
	}

	if timesExceeded >= 3 {
		budgetDetailDTO.ExceededColor = "danger"
	} else if timesExceeded >= 1 {
		budgetDetailDTO.ExceededColor = "warning"
	} else {
		budgetDetailDTO.ExceededColor = "success"
	}

	if float64(spendStdDeviation)/float64(averageSpend) > 0.5 {
		budgetDetailDTO.Volatility = "Very High"
	} else if float64(spendStdDeviation)/float64(averageSpend) > 0.25 {
		budgetDetailDTO.Volatility = "High"
	} else {
		budgetDetailDTO.Volatility = "Low"
	}

	tmpl, err := template.New("budgetDetail").Parse(budgetDetailTmpl)
	if err != nil {
		message := fmt.Sprintf("Error building template: %v", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	err = utils.RenderTemplateAsHTML(w, tmpl, budgetDetailDTO)
	if err != nil {
		panic(err)
	}

}
