package services

import (
	"encoding/csv"
	"strings"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type Parser interface {
	Parse(string) ([]models.Transaction, []models.Balance, error)
}

type SchwabCheckingCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column,
// description in 4th column, withdrawal amount in 5th column,
// deposit Amount in 6th column, and running balance in 7th column
func (s SchwabCheckingCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[0])
		if idx == 1 {
			balance := utils.DollarStringToCents(record[7])
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		if record[5] != "" {
			amount = utils.DollarStringToCents(record[5])
		} else if record[6] != "" {
			amount = utils.DollarStringToCents(record[6])
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[4],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, balances, nil
}

type SchwabBrokerageCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column,
// action in 1st column, symbol in 2nd column,
// description in 3rd column, amount in 7th column,
func (s SchwabBrokerageCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		// Schwab Brokerage reports sometimes include a date value like
		// "09/30/2024 as of 09/29/2024" so we need to extract the date
		date := strings.Split(record[0], " ")[0]

		isoDate := utils.ConvertMMDDYYYYtoISO8601(date)
		var amount int64
		amount = utils.DollarStringToCents(record[7])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[1] + " - " + record[3],
			Amount:      amount,
		}
		transactions = append(transactions, txn)
	}
	return transactions, balances, nil
}

type FidelityCreditCardCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
func (FidelityCreditCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[0])
		var amount int64
		amount = utils.DollarStringToCents(record[4])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, []models.Balance{}, nil
}

type FidelityBrokerageCSVParser struct{}

// Parses CSVs with the header as the 2nd row, date in 0th column,
// description in 1st column, amount in 10th column, and balance in 11th column
// Transactions are sorted by newest transaction first, so the balance is the
// first row after the header
func (FidelityBrokerageCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	csvReader := csv.NewReader(strings.NewReader(statement))
	// Fidelity includes extra disclosures at the end of their brokerage CSVs
	// so we need to disable FieldsPerRecord column count validation
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header rows
		if idx < 3 {
			continue
		}
		// Fidelity includes extra disclosures at the end of their brokerage
		// CSVs so we drop any records that don't have all columns
		if len(record) < 13 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[0])
		if idx == 3 {
			balance := utils.DollarStringToCents(record[11])
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		amount = utils.DollarStringToCents(record[10])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[1],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, balances, nil
}

type ChaseCreditCardCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
func (s ChaseCreditCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[0])
		var amount int64
		amount = utils.DollarStringToCents(record[4])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, []models.Balance{}, nil
}

type CapitalOneCredictCardCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 3rd column, category in 4th column, debits (purchases) in 5th column,
// credit (payments/refunds) amount in 6th column
func (s CapitalOneCredictCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[0])
		var amount int64
		if record[5] != "" {
			amount = utils.DollarStringToCents(record[5])
			// Negate the amount since it's a debit
			amount = amount * -1
		} else if record[6] != "" {
			amount = utils.DollarStringToCents(record[6])
		}
		// TODO: use category from capital one to set category in transaction
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[3],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, []models.Balance{}, nil
}

type CapitalOneSavingsCSVParser struct{}

// Parses CSVs with the header as the 1st row,  description in 1st column, date
// in 2nd column, transaction type (credit vs debit) in 3rd column, amount in
// 4th column, and balance in 5th column. Transactions are sorted by newest
// transaction first, so the balance is the first row after the header
func (s CapitalOneSavingsCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[2])
		if idx == 1 {
			balance := utils.DollarStringToCents(record[5])
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		amount = utils.DollarStringToCents(record[4])

		if record[3] == "Debit" {
			amount = amount * -1
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[1],
			Amount:      amount,
		}
		transactions = append(transactions, txn)
	}
	return transactions, balances, nil
}

var parsersByInstitution map[string]Parser = map[string]Parser{
	"schwabChecking":       SchwabCheckingCSVParser{},
	"schwabBrokerage":      SchwabBrokerageCSVParser{},
	"fidelityCreditCard":   FidelityCreditCardCSVParser{},
	"fidelityBrokerage":    FidelityBrokerageCSVParser{},
	"chaseCreditCard":      ChaseCreditCardCSVParser{},
	"capitalOneCreditCard": CapitalOneCredictCardCSVParser{},
	"capitalOneSavings":    CapitalOneSavingsCSVParser{},
}
