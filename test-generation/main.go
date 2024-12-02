package main

import (
	"fmt"
	"io"
	"text/template"
)

var testFuncString = `func test_{{ csvName }}Parser() error {
	parser := parsersByInstitution["bankOfAmericaCreditCard"]
	parser.Parse({{ csvName }})
	txns, balances, err := parser.Parse({{ csvName }})
	if err != nil {
		fmt.Println("failed to parse: ", err)
	}
	if len(txns) != 6 {
		fmt.Printf("transaction count is wrong - got %s, expected %s\n", 6, len(txns))
	}
	if len(balances) != 6 {
		fmt.Printf("balances count is wrong - got %s, expected %s\n", 6, len(txns))
	}

	return nil
}`

func main() {
	tmpl, _ := template.New("whatever").Parse(testFuncString)
	tmplData := map[string]string{csvName: "bankOfAmericaCreditCard"}
	foo := io.Writer{}
	tmpl.Execute(nil, tmplData)
	fmt.Println("test")
}
