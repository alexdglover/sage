package services

import "fmt"

var happyPathCSV = ""

// For testing cases where there is a footer in the CSV, or inconsistent record length
var csvWithExtraShit = ""

// Parses CSVs with the header as the 1st row, date in 0th column, description in 2nd column, and amount in 4th column
var bankOfAmericaCreditCardCSV = ""

// Parses CSVs with the header as the 1st row, date in 1st column, description
// in 2nd column, and amount in 3rd column
var chaseCheckingCSV = ""

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
var chaseCreditCardCSV = ""

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 3rd column, category in 4th column, debits (purchases) in 5th column,
// credit (payments/refunds) amount in 6th column
var capitalOneCreditCardCSV = ""

// Parses CSVs with the header as the 1st row,  description in 1st column, date
// in 2nd column, transaction type (credit vs debit) in 3rd column, amount in
// 4th column, and balance in 5th column. Transactions are sorted by newest
// transaction first, so the balance is the first row after the header
var capitalOneSavingsCSV = ""

// Parses CSVs with the header as the 1st row, date in 0th column,
// description in 4th column, withdrawal amount in 5th column,
// deposit Amount in 6th column, and running balance in 7th column
var schwabCheckingCSV = ""

// Parses CSVs with the header as the 1st row, date in 0th column,
// action in 1st column, symbol in 2nd column,
// description in 3rd column, amount in 7th column,
var schwabBrokerageCSV = ""

// Parses CSVs with the header as the 1st row, date in 0th column, description
// in 2nd column, and amount in 4th column
var fidelityCreditCardCSV = ""

// Parses CSVs with the header as the 2nd row, date in 0th column,
// description in 1st column, amount in 10th column, and balance in 11th column
// Transactions are sorted by newest transaction first, so the balance is the
// first row after the header
var fidelityBrokerageCSV = ""

func test_generalCsvParser() error {
	generalCSVParser := GeneralCSVParser{}
	txns, balances, err := generalCSVParser.Parse(happyPathCSV, 0, 1, 2, true, false)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_BankOfAmericaCreditCardCSVParser() error {
	parser := parsersByInstitution["bankOfAmericaCreditCard"]
	txns, balances, err := parser.Parse(bankOfAmericaCreditCardCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_SchwabCheckingCSVParser() error {
	parser := parsersByInstitution["schwabChecking"]
	txns, balances, err := parser.Parse(schwabCheckingCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_SchwabBrokerageCSVParser() error {
	parser := parsersByInstitution["schwabBrokerage"]
	txns, balances, err := parser.Parse(schwabBrokerageCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_FidelityCreditCardCSVParser() error {
	parser := parsersByInstitution["fidelityCreditCard"]
	txns, balances, err := parser.Parse(fidelityCreditCardCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_FidelityBrokerageCSVParser() error {
	parser := parsersByInstitution["fidelityBrokerage"]
	txns, balances, err := parser.Parse(fidelityBrokerageCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_ChaseCreditCardCSVParser() error {
	parser := parsersByInstitution["chaseCreditCard"]
	txns, balances, err := parser.Parse(chaseCreditCardCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_ChaseCheckingCSVParser() error {
	parser := parsersByInstitution["chaseChecking"]
	txns, balances, err := parser.Parse(chaseCheckingCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_CapitalOneCreditCardCSVParser() error {
	parser := parsersByInstitution["capitalOneCreditCard"]
	txns, balances, err := parser.Parse(capitalOneCreditCardCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}

func test_CapitalOneSavingsCSVParser() error {
	parser := parsersByInstitution["capitalOneSavings"]
	txns, balances, err := parser.Parse(capitalOneSavingsCSV)
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %v, expected %v\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %v, expected %v\n", 6, len(txns))
	}

	return nil
}
