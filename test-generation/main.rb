

def generate_test_string(institution_name)
    test_name = institution_name.clone()
    first_char = test_name[0]
    first_char.capitalize!
    test_name[0] = first_char
    puts "func test_#{ test_name }CSVParser() error {
    parser := parsersByInstitution[\"#{ institution_name }\"]
    txns, balances, err := parser.Parse(#{ institution_name }CSV)
    if err != nil {
        fmt.Println(\"failed to parse: \", err)
    }
    if len(txns) != 6 {
        fmt.Printf(\"transaction count is wrong - got %s, expected %s\\n\", 6, len(txns))
    }
    if len(balances) != 6 {
        fmt.Printf(\"balances count is wrong - got %s, expected %s\\n\", 6, len(txns))
    }

    return nil
}\n\n"
end

["bankOfAmericaCreditCard", "schwabChecking", "schwabBrokerage", "fidelityCreditCard", "fidelityBrokerage", "chaseCreditCard", "chaseChecking", "capitalOneCreditCard", "capitalOneSavings"].each do |institution|
    generate_test_string(institution)
end