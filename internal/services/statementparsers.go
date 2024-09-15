package services

import (
	"encoding/csv"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexdglover/sage/internal/models"
)

type Parser interface {
	Parse(string) ([]models.Transaction, []models.Balance, error)
}

type SchwabCSVParser struct{}

// Parses CSVs with the header as the 1st row, Date in col0, Description in
// col4, withdrawal Amount col5, deposit Amount in col6, and running
// balance in col7
func (s SchwabCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println("records are", records)
	for idx, record := range records {
		// fmt.Println("working on record", idx)
		// Skip the header row
		if idx == 0 {
			continue
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
		// fmt.Println("amount is", amount)
		txn := models.Transaction{
			Date:        record[0],
			Description: record[4],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
	return transactions, []models.Balance{}, nil
}

type FidelityCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
func (FidelityCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println("records are", records)
	for idx, record := range records {
		// fmt.Println("working on record", idx)
		// Skip the header row
		if idx == 0 {
			continue
		}
		var amount int64
		amount, err = parseAmount(record[4])
		if err != nil {
			return nil, nil, err
		}
		// fmt.Println("amount is", amount)
		txn := models.Transaction{
			Date:        record[0],
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
	return transactions, []models.Balance{}, nil
}

type ChaseCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
func (s ChaseCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println("records are", records)
	for idx, record := range records {
		// fmt.Println("working on record", idx)
		// Skip the header row
		if idx == 0 {
			continue
		}
		var amount int64
		amount, err = parseAmount(record[4])
		if err != nil {
			return nil, nil, err
		}
		// fmt.Println("amount is", amount)
		txn := models.Transaction{
			Date:        record[0],
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
		// fmt.Println("amount is", amount)
		txn := models.Transaction{
			Date:        record[0],
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
		if idx == 1 {
			balance, err := parseAmount(record[5])
			if err != nil {
				return nil, nil, err
			}
			balances = append(balances, models.Balance{
				EffectiveDate: record[2],
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
		// fmt.Println("amount is", amount)
		txn := models.Transaction{
			Date:        record[2],
			Description: record[1],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	// for each row in the CSV, parse the columns and add it to transactions
	return transactions, []models.Balance{}, nil
}

var parsersByInstitution map[string]Parser = map[string]Parser{
	"schwab":               SchwabCSVParser{},
	"fidelity":             FidelityCSVParser{},
	"chase":                ChaseCSVParser{},
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
