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

type GeneralCSVParser struct{}

var generalCSVParser = GeneralCSVParser{}

func (g GeneralCSVParser) Parse(statement string, dateCol int, descCol int, amountCol int, skipHeader bool, skipRecordLengthValidation bool) (transactions []models.Transaction, balances []models.Balance, err error) {
	// parse the string into a CSV
	csvReader := csv.NewReader(strings.NewReader(statement))
	// Some exports include a trailing comma or extra commentary in the data
	// In those cases we need to disable FieldsPerRecord column count validation
	if skipRecordLengthValidation {
		csvReader.FieldsPerRecord = -1
	}
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 && skipHeader {
			continue
		}
		isoDate := utils.ConvertMMDDYYYYtoISO8601(record[dateCol])
		amount := utils.DollarStringToCents(record[amountCol])
		if amount < 0 {
			amount = amount * -1
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[descCol],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, []models.Balance{}, nil
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
		var amount int
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
		var amount int
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
		amount := utils.DollarStringToCents(record[4])
		// We negate by transaction category rather than at the amount property, so convert any negative amounts into positive
		if amount < 0 {
			amount = amount * -1
		}
		txn := models.Transaction{
			Date:        record[0],
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
		var amount int
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

type ChaseCheckingCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 1st column, description
// in 2nd column, and amount in 3rd column
func (ChaseCheckingCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	return generalCSVParser.Parse(statement, 1, 2, 3, true, true)
}

type BankOfAmericaCreditCardCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description in 2nd column, and amount in 4th column
func (BankOfAmericaCreditCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	return generalCSVParser.Parse(statement, 0, 2, 4, true, false)
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
		amount := utils.DollarStringToCents(record[5])
		if amount < 0 {
			amount = amount * -1
		}
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[2],
			Amount:      amount,
		}
		transactions = append(transactions, txn)

	}
	return transactions, []models.Balance{}, nil
}

type CapitalOneCreditCardCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 3rd column, category in 4th column, debits (purchases) in 5th column,
// credit (payments/refunds) amount in 6th column
func (s CapitalOneCreditCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
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
		var amount int
		if record[5] != "" {
			amount = utils.DollarStringToCents(record[5])
		} else if record[6] != "" {
			amount = utils.DollarStringToCents(record[6])
		}
		// TODO: use category from capital one to set category in transaction
		txn := models.Transaction{
			Date:        record[0],
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
		isoDate := utils.ConvertMMDDYYtoISO8601(record[2])
		if idx == 1 {
			balance := utils.DollarStringToCents(record[5])
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		amount := utils.DollarStringToCents(record[4])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[1],
			Amount:      amount,
		}
		transactions = append(transactions, txn)
	}
	return transactions, balances, nil
}

type TargetCreditCardCSVParser struct{}

// Parses CSVs with the header as the 1st row,  date in 0th column, posting date
// in 1st column, ref# in 2nd column (which we won't use), amount in 3rd column,
// description in 4th column, last 4 digits of card number in 5th column, and transaction
// type in 6th column. Transaction typ[e is either `Payment`, `Sale`, or `Refund`.
func (s TargetCreditCardCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	csvReader := csv.NewReader(strings.NewReader(statement))
	// Target statement CSVs include empty fields with double quotes,
	// which is interpreted as an escaped double quote to the parser.
	// To disable this behavior, we need to set the LazyQuotes flag to true.
	csvReader.LazyQuotes = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		amount := utils.DollarStringToCents(record[3])
		txn := models.Transaction{
			Date:        record[0],
			Description: record[4],
			Amount:      amount,
		}
		transactions = append(transactions, txn)
	}
	return transactions, balances, nil
}

type UWCUMortgageCSVParser struct{}

// Parses CSVs with the header as the 1st row, date in 2nd column, amount in 3rd column,
// description in 4th column, and balance in the 7th column. Transactions are sorted by newest
// transaction first, so the balance is the first row after the header
func (s UWCUMortgageCSVParser) Parse(statement string) (transactions []models.Transaction, balances []models.Balance, err error) {
	csvReader := csv.NewReader(strings.NewReader(statement))
	// Statement CSVs include empty fields with double quotes,
	// which is interpreted as an escaped double quote to the parser.
	// To disable this behavior, we need to set the LazyQuotes flag to true.
	csvReader.LazyQuotes = true
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	for idx, record := range records {
		// Skip the header row
		if idx == 0 {
			continue
		}
		isoDate := utils.ConvertMDYYYYtoISO8601(record[2])
		if idx == 1 {
			balance := utils.DollarStringToCents(record[7])
			balances = append(balances, models.Balance{
				EffectiveDate: isoDate,
				Amount:        balance,
			})
		}
		amount := utils.DollarStringToCents(record[3])
		txn := models.Transaction{
			Date:        isoDate,
			Description: record[4],
			Amount:      amount,
		}
		transactions = append(transactions, txn)
	}
	return transactions, balances, nil
}

var parsersByInstitution map[string]Parser = map[string]Parser{
	"bankOfAmericaCreditCard": BankOfAmericaCreditCardCSVParser{},
	"capitalOneCreditCard":    CapitalOneCreditCardCSVParser{},
	"capitalOneSavings":       CapitalOneSavingsCSVParser{},
	"chaseCreditCard":         ChaseCreditCardCSVParser{},
	"chaseChecking":           ChaseCheckingCSVParser{},
	"fidelityBrokerage":       FidelityBrokerageCSVParser{},
	"fidelityCreditCard":      FidelityCreditCardCSVParser{},
	"schwabChecking":          SchwabCheckingCSVParser{},
	"schwabBrokerage":         SchwabBrokerageCSVParser{},
	"targetCreditCard":        TargetCreditCardCSVParser{},
	"uwcuMortgage":            UWCUMortgageCSVParser{},
}
