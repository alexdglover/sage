package services

import (
	"encoding/csv"
	"regexp"
	"strconv"
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
			balance, err := parseAmount(record[7])
			if err != nil {
				return nil, nil, err
			}
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		if record[5] != "" {
			amount, err = parseAmount(record[5])
			if err != nil {
				return nil, nil, err
			}
		} else if record[6] != "" {
			amount, err = parseAmount(record[6])
			if err != nil {
				return nil, nil, err
			}
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[4],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
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
		amount, err = parseAmount(record[4])
		if err != nil {
			return nil, nil, err
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
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
			balance, err := parseAmount(record[11])
			if err != nil {
				return nil, nil, err
			}
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		amount, err = parseAmount(record[10])
		if err != nil {
			return nil, nil, err
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[1],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
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
		amount, err = parseAmount(record[4])
		if err != nil {
			return nil, nil, err
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
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
			amount, err = parseAmount(record[5])
			// Negate the amount since it's a debit
			amount = amount * -1
			if err != nil {
				return nil, nil, err
			}
		} else if record[6] != "" {
			amount, err = parseAmount(record[6])
			if err != nil {
				return nil, nil, err
			}
		}
		// TODO: use category from capital one to set category in transaction
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[3],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
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
			balance, err := parseAmount(record[5])
			if err != nil {
				return nil, nil, err
			}
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		var amount int64
		amount, err = parseAmount(record[4])
		if err != nil {
			return nil, nil, err
		}

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
	"fidelityCreditCard":   FidelityCreditCardCSVParser{},
	"fidelityBrokerage":    FidelityBrokerageCSVParser{},
	"chaseCreditCard":      ChaseCreditCardCSVParser{},
	"capitalOneCreditCard": CapitalOneCredictCardCSVParser{},
	"capitalOneSavings":    CapitalOneSavingsCSVParser{},
}

func parseAmount(amount string) (amountAsInt int64, err error) {
	var amountAsFloat float64
	re := regexp.MustCompile(`[^0-9.-]`)
	amount = re.ReplaceAllString(amount, "")
	amountAsFloat, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}
	amountAsInt = int64(amountAsFloat * 100)
	return amountAsInt, nil
}
