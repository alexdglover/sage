package services

import (
	"fmt"

	"github.com/GopherML/bag"
	"github.com/alexdglover/sage/internal/models"
)

type MLCategorizer struct {
	Bag                   *bag.Bag
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

func (mc *MLCategorizer) CategorizeTransaction(transaction *models.Transaction) (string, error) {
	results := mc.Bag.GetResults(transaction.Description)
	category := results.GetHighestProbability()
	fmt.Println("Categorizing transaction: ", transaction.Description, " as ", category, " with score ", results[category])
	return category, nil
}
