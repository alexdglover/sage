package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type NetWorthController struct {
	BalanceRepository *models.BalanceRepository
}

//go:embed networth.html
var netWorthTmpl string

// A totalByType entry will be either 'asset': <some int> or 'liablities': <some int>
// This is used for aggregating balances, but is not the DTO to the front end HTML template
type totalByType map[string]int

// The DTO is also a map of year-month to amount, but has the total converted to string for
// expected currency format in the UI
type totalByTypeDTO map[string]string

type netWorthdto struct {
	ActivePage          string
	AllTimeActive       bool
	Last12MonthsActive  bool
	Last6MonthsActive   bool
	Last3MonthsActive   bool
	TotalByMonthAndType map[string]totalByTypeDTO
}

// TODO: Consider moving this into a service class that returns just the data needed
func (nc *NetWorthController) netWorthHandler(w http.ResponseWriter, req *http.Request) {
	var relativeWindow int
	dto := netWorthdto{
		ActivePage: "netWorth",
	}
	if req.FormValue("relativeWindow") == "" {
		relativeWindow = 6
		dto.Last6MonthsActive = true
	} else {
		switch req.FormValue("relativeWindow") {
		case "12":
			relativeWindow = 12
			dto.Last12MonthsActive = true
		case "6":
			relativeWindow = 6
			dto.Last6MonthsActive = true
		case "3":
			relativeWindow = 3
			dto.Last3MonthsActive = true
		case "allTime":
			// 10 years in months, as a useful approximation of all time for now
			relativeWindow = 120
			dto.AllTimeActive = true
		default:
			fmt.Println("invalid relative window provided, falling back to 6 months")
			relativeWindow = 6
			dto.Last6MonthsActive = true
		}
	}

	// We always start with today's date and work backwards based on relative window value
	endDate := time.Now()
	// Decrement relativeWindow by 1 (to account for the current month already being included)
	relativeWindow = relativeWindow - 1
	// And calculate start date
	startDate := endDate.AddDate(0, (relativeWindow * -1), 0)

	assetBalances := nc.BalanceRepository.GetBalancesOfAllAssetsByMonth(context.TODO(), startDate, endDate)
	liabilityBalances := nc.BalanceRepository.GetBalancesOfAllLiabilitiesByMonth(context.TODO(), startDate, endDate)

	// A map of year-month to total by type, so we can aggregate the sum of balances for each year-month
	totalByMonthAndType := map[string]totalByType{}
	// A map of year-month to total by type as a string value, so it can be presented in the UI layer
	// this keeps the aggregation math separate from presentation layer
	totalByMonthAndTypeDTO := map[string]totalByTypeDTO{}

	for _, balancesByDate := range assetBalances {
		yearMonth := balancesByDate.Date.Format("2006-01")
		if totalByMonthAndType[yearMonth] == nil {
			totalByMonthAndType[yearMonth] = totalByType{}
		}

		for _, balance := range balancesByDate.Balances {
			totalByMonthAndType[yearMonth]["assets"] = totalByMonthAndType[yearMonth]["assets"] + balance.Amount
		}
	}

	for _, balancesByDate := range liabilityBalances {
		yearMonth := balancesByDate.Date.Format("2006-01")
		if totalByMonthAndType[yearMonth] == nil {
			totalByMonthAndType[yearMonth] = totalByType{}
		}
		for _, balance := range balancesByDate.Balances {
			totalByMonthAndType[yearMonth]["liabilities"] = totalByMonthAndType[yearMonth]["liabilities"] - balance.Amount
		}
	}

	for date := range totalByMonthAndType {
		totalByMonthAndType[date]["netWorth"] = totalByMonthAndType[date]["assets"] + totalByMonthAndType[date]["liabilities"]
		// populate the totalByTypeDTO
		if totalByMonthAndTypeDTO[date] == nil {
			totalByMonthAndTypeDTO[date] = totalByTypeDTO{}
		}
		totalByMonthAndTypeDTO[date]["assets"] = utils.CentsToDollarStringMachineSafe(totalByMonthAndType[date]["assets"])
		totalByMonthAndTypeDTO[date]["liabilities"] = utils.CentsToDollarStringMachineSafe(totalByMonthAndType[date]["liabilities"])
		totalByMonthAndTypeDTO[date]["netWorth"] = utils.CentsToDollarStringMachineSafe(totalByMonthAndType[date]["netWorth"])
		// humanized versions for display
		totalByMonthAndTypeDTO[date]["humanizedLiabilities"] = utils.CentsToDollarStringHumanized(totalByMonthAndType[date]["liabilities"])
		totalByMonthAndTypeDTO[date]["humanizedAssets"] = utils.CentsToDollarStringHumanized(totalByMonthAndType[date]["assets"])
		totalByMonthAndTypeDTO[date]["humanizedNetWorth"] = utils.CentsToDollarStringHumanized(totalByMonthAndType[date]["netWorth"])

	}

	dto.TotalByMonthAndType = totalByMonthAndTypeDTO
	tmpl := template.Must(template.New("netWorthDashboard").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(netWorthTmpl))
	err := utils.RenderTemplateAsHTML(w, tmpl, dto)
	if err != nil {
		panic(err)
	}
}
