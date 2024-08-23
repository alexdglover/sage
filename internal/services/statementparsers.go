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

// Parses CSVs with the header as the 1st row, Date in col0, Description in col4, withdrawal Amount col5,
// deposit Amount in col6, and running Balance in col7
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

type FidelityVisaCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description in 2nd column, and amount in 4th column
func (FidelityVisaCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
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

type ChaseVisaCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description in 2nd column, and amount in 4th column
func (s ChaseVisaCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
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

var parsersByInstitution map[string]Parser = map[string]Parser{
	"schwab": SchwabCSVParser{},
	"fidelity visa": FidelityVisaCSVParser{},

}

func parseAmount(amount string) (amountAsInt int64, err error) {
	var amountAsFloat float64
	re := regexp.MustCompile(`[^0-9.]`)
	amount = re.ReplaceAllString(amount, "")
	amountAsFloat, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}
	amountAsInt = int64(amountAsFloat * 100)
	return amountAsInt, nil
}
