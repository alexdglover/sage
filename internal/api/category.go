package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type CategoryController struct {
	CategoryRepository *models.CategoryRepository
	BalanceRepository  *models.BalanceRepository
}

//go:embed categories.html
var categoriesPageTmpl string

//go:embed categoryForm.html
var categoryFormTmpl string

type CategoryDTO struct {
	ID   uint
	Name string
}

type CategoriesPageDTO struct {
	Categories          []CategoryDTO
	CategorySaved       bool
	CreatedCategoryName string
}

type CategoryFormDTO struct {
	// If we're updating an existing category in the form, Updating will be true
	// If we're creating a new category, Updating will be false
	Updating     bool
	CategoryID   string
	CategoryName string
}

func (ac *CategoryController) generateCategoriesView(w http.ResponseWriter, req *http.Request) {
	// Get all categories
	categories, err := ac.CategoryRepository.GetAllCategories()
	if err != nil {
		http.Error(w, "Unable to get categories", http.StatusInternalServerError)
		return
	}

	// Build categories DTO
	categoriesDTO := []CategoryDTO{}
	for _, category := range categories {
		// Skip the "Unknown" category because we don't want the user to edit/delete it
		if category.Name == "Unknown" {
			continue
		}
		categoriesDTO = append(categoriesDTO, CategoryDTO{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	categoriesPageDTO := CategoriesPageDTO{
		Categories: categoriesDTO,
	}
	if req.URL.Query().Get("categorySaved") != "" {
		categoriesPageDTO.CategorySaved = true
		categoriesPageDTO.CreatedCategoryName = req.URL.Query().Get("categorySaved")
	}

	tmpl := template.Must(template.New("categoriesPage").Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
	}).Parse(categoriesPageTmpl))

	err = tmpl.Execute(w, categoriesPageDTO)
	if err != nil {
		panic(err)
	}
}

func (ac *CategoryController) generateCategoryForm(w http.ResponseWriter, req *http.Request) {
	var dto CategoryFormDTO

	categoryIDQueryParameter := req.URL.Query().Get("categoryID")
	if categoryIDQueryParameter != "" {
		categoryID, err := utils.StringToUint(categoryIDQueryParameter)
		if err != nil {
			http.Error(w, "Unable to parse category ID", http.StatusInternalServerError)
			return
		}
		category, err := ac.CategoryRepository.GetCategoryByID(categoryID)
		if err != nil {
			http.Error(w, "Unable to get category", http.StatusInternalServerError)
			return
		}

		dto = CategoryFormDTO{
			Updating:     true,
			CategoryID:   fmt.Sprint(category.ID),
			CategoryName: category.Name,
		}
	}

	tmpl, err := template.New("categoryForm").Parse(categoryFormTmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, dto)
	if err != nil {
		panic(err)
	}
}

func (ac *CategoryController) upsertCategory(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	categoryID := req.FormValue("categoryID")
	categoryName := req.FormValue("categoryName")

	var category models.Category

	if categoryID != "" {
		id, err := utils.StringToUint(categoryID)
		if err != nil {
			http.Error(w, "Unable to parse category ID", http.StatusBadRequest)
			return
		}
		category, err = ac.CategoryRepository.GetCategoryByID(id)
		if err != nil {
			http.Error(w, "Unable to get category", http.StatusBadRequest)
			return
		}
	} else {
		category = models.Category{}
	}

	category.Name = categoryName

	_, err := ac.CategoryRepository.Save(category)
	if err != nil {
		http.Error(w, "Unable to save category", http.StatusBadRequest)
		return
	}

	queryValues := url.Values{}
	queryValues.Add("categorySaved", categoryName)
	// TODO: Consider moving the categoryView to a function that accepts an extra argument
	// instead of invoking the endpoint with a custom request
	categoryViewReq := http.Request{
		Method: "GET",
		URL: &url.URL{
			RawQuery: queryValues.Encode(),
		},
	}
	categoryViewReq.URL.RawQuery = queryValues.Encode()

	ac.generateCategoriesView(w, &categoryViewReq)
}
