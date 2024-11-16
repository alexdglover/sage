package services

import (
	"fmt"

	"github.com/GopherML/bag"
	"github.com/alexdglover/sage/internal/models"
)

type MLCategorizer struct {
	Bag                   *bag.Bag
	CategoryRepository    *models.CategoryRepository
	TransactionRepository *models.TransactionRepository
}

func (mc *MLCategorizer) BuildModel() error {
	// Get all transactions flagged for training
	transactions, err := mc.TransactionRepository.GetTransactionsForTraining()
	if err != nil {
		return err
	}
	bagConfig := bag.Config{}
	mc.Bag, err = bag.New(bagConfig)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		mc.Bag.Train(transaction.Description, transaction.Category.Name)
	}

	return nil
}

// Should this return just the category name as a string, a category object, both, or something else?
func (mc *MLCategorizer) CategorizeTransaction(transaction *models.Transaction) (category models.Category, err error) {
	fmt.Println("yo categorizetransaction was invoked")
	fmt.Println("bag is nil? ", mc.Bag == nil)
	results := mc.Bag.GetResults(transaction.Description)
	fmt.Println("Results: ", results)
	categoryName := results.GetHighestProbability()
	category, err = mc.CategoryRepository.GetCategoryByName(categoryName)
	if err != nil {
		category = models.Category{}
		return category, err

	}
	fmt.Println("Categorizing transaction: ", transaction.Description, " as ", category, " with score ", results[categoryName])
	return category, nil
}
