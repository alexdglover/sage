package utils

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
)

func StringToUint(input string) (uint, error) {
	// Convert string to uint first
	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, err
	}
	// Convert uint to uint and return
	return uint(output), nil
}

// Convert whole cents (as int) to dollars (as float64) and returns a "commaized" string
func CentsToDollarStringHumanized(input int) string {
	amount := float64(input) / 100
	return trimTrailingDigits(humanize.Commaf(amount), 2)
}

// Convert whole cents (as int) to dollars (as float64) and returns a string
func CentsToDollarStringMachineSafe(input int) string {
	amount := float64(input) / 100
	return fmt.Sprintf("%.2f", amount)
}

// Trims digits following decimal place in "number strings"
func trimTrailingDigits(s string, digits int) string {
	if digits <= 0 {
		return s
	}
	if i := strings.Index(s, "."); i >= 0 {
		if len(s[i+1:]) < digits {
			return s + strings.Repeat("0", digits-len(s[i+1:]))
		}
		if digits <= 0 {
			return s[:i]
		}
		i++
		if i+digits >= len(s) {
			return s
		}
		return s[:i+digits]
	}
	return s + "." + strings.Repeat("0", digits)
}

func DollarStringToCents(input string) int {
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
	amountAsInt := int(amountAsFloat * 100)
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
	sanitizedInput := strings.TrimSpace(input)
	t, err := time.Parse("01/02/2006", sanitizedInput)
	if err != nil {
		panic(err)
	}
	return TimeToISO8601DateString(t)
}

func ConvertMMDDYYtoISO8601(input string) string {
	sanitizedInput := strings.TrimSpace(input)
	t, err := time.Parse("01/02/06", sanitizedInput)
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

func Percentile(input []int, percentile float64) int {
	if len(input) == 0 {
		return 0
	}
	sortedDataSet := make([]int, len(input))
	copy(sortedDataSet, input)
	sort.Ints(sortedDataSet)

	for idx, value := range sortedDataSet {
		// If we're at the index that corresponds to the percentile we're looking for, return the value
		percentileIndexAsFloat := percentile * float64(len(sortedDataSet))
		percentileIndex := int(math.Ceil(percentileIndexAsFloat))
		if (idx + 1) == percentileIndex {
			return value
		}
	}
	// This shouldn't be reachable, but needed to make compiler happy
	return 0
}

func RenderTemplateAsHTML(w http.ResponseWriter, tmpl *template.Template, data any) (err error) {
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	return err
}
