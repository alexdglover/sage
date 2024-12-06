package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func StringToUint(input string) (uint, error) {
	// Convert string to uint64 first
	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, err
	}
	// Convert uint64 to uint and return
	return uint(output), nil
}

func CentsToDollarString(input int64) string {
	// Convert whole cents (as int64) to dollars (as float64)
	amount := float64(input) / 100
	// Convert float to string with 2 decimal places, to force 2 decimal places
	return fmt.Sprintf("%.2f", amount)
}

func DollarStringToCents(input string) int64 {
	if input == "" {
		return 0
	}

	// Remove all non-numeric characters (like commas, dollar signs) from the input string
	re := regexp.MustCompile(`[^0-9.-]`)
	amount := re.ReplaceAllString(input, "")

	amountAsFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	amountAsInt := int64(amountAsFloat * 100)
	return amountAsInt
}

func TimeToISO8601DateString(input time.Time) string {
	return fmt.Sprint(input.Format("2006-01-02"))
}

func ISO8601DateStringToTime(input string) time.Time {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		panic(err)
	}
	return t
}

func ConvertMMDDYYYYtoISO8601(input string) string {
	t, err := time.Parse("01/02/2006", input)
	if err != nil {
		panic(err)
	}
	return TimeToISO8601DateString(t)
}

func ConvertMMDDYYtoISO8601(input string) string {
	t, err := time.Parse("01/02/06", input)
	if err != nil {
		panic(err)
	}
	return TimeToISO8601DateString(t)
}

func StrPointer(input string) *string {
	return &input
}

func StrPointerToString(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func UintPointerToString(input *uint) string {
	if input == nil {
		return ""
	}
	return fmt.Sprint(*input)
}

func AmountValid(input string) bool {
	if input == "" {
		return false
	}
	var validAmount = regexp.MustCompile(`^[0-9]*[\.]{0,1}[0-9]{0,2}`)
	return validAmount.MatchString(input)
}

// DateValid returns true if the input string is a valid ISO8601 date, specifically in the YYYY-MM-DD format
func DateValid(input string) bool {
	if input == "" {
		return false
	}
	var validIso8601Date = regexp.MustCompile(`^\d{4}-([0][1-9]|1[0-2])-([0][1-9]|[1-2]\d|3[01])$`)
	return validIso8601Date.MatchString(input)
}
