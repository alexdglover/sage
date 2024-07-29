package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

//go:embed networth.html.tmpl
var netWorthTmpl string

// A dataPoint will be either 'asset': <some float> or 'liablities': <some float>
type DataPoint map[string]float32

type netWorthTmplVariables struct {
	AllTimeActive      bool
	Last12MonthsActive bool
	Last6MonthsActive  bool
	Last3MonthsActive  bool
	DataPoints         map[string]DataPoint
}

// TODO: Consider moving this into a service class that returns just the data needed
func netWorthHandler(w http.ResponseWriter, req *http.Request) {
	var relativeWindow int
	tmplVariables := netWorthTmplVariables{}
	if req.FormValue("relativeWindow") == "" {
		relativeWindow = 6
		tmplVariables.Last6MonthsActive = true
	} else {
		switch req.FormValue("relativeWindow") {
		case "12":
			relativeWindow = 12
			tmplVariables.Last12MonthsActive = true
		case "6":
			relativeWindow = 6
			tmplVariables.Last6MonthsActive = true
		case "3":
			relativeWindow = 3
			tmplVariables.Last3MonthsActive = true
		case "allTime":
			// 100 years in months, as a useful approximation of all time
			relativeWindow = 1200
			tmplVariables.AllTimeActive = true
		default:
			fmt.Println("invalid relative window provided, falling back to 6 months")
			relativeWindow = 6
			tmplVariables.Last6MonthsActive = true
		}
	}

	// We always start with today's date and work backwards based on relative window value
	endDate := time.Now()
	// Decrement relativeWindow by 1 (to account for current month already being included)
	relativeWindow = relativeWindow - 1
	// And calculate start date
	startDate := endDate.AddDate(0, (relativeWindow * -1), 0)
	// Iterate over each month between start and end dates to build dateRange slice
	// for date := startDate; !date.After(endDate); date = date.AddDate(0, 1, 0) {
	// 	dateRange = append(dateRange, date)
	// }

	br := models.GetBalanceRepository()
	// assetBalances := br.GetBalancesOfAllAssets(context.TODO(), startDate, endDate)
	assetBalances := br.GetBalancesOfAllAssetsByMonth(context.TODO(), startDate, endDate)
	// fmt.Println("asset balances are: ", assetBalances)

	// liabilityBalances := br.GetBalancesOfAllLiabilities(context.TODO(), startDate, endDate)
	liabilityBalances := br.GetBalancesOfAllLiabilitiesByMonth(context.TODO(), startDate, endDate)
	// fmt.Println("liability balances are: ", liabilityBalances)

	dataPointSet := map[string]DataPoint{}

	for _, balancesByDate := range assetBalances {
		yearMonth := balancesByDate.Date.Format("2006-01")
		if dataPointSet[yearMonth] == nil {
			dataPointSet[yearMonth] = DataPoint{}
		}

		for _, balance := range balancesByDate.Balances {
			dataPointSet[yearMonth]["assets"] = dataPointSet[yearMonth]["assets"] + balance.Balance
		}
	}

	for _, balancesByDate := range liabilityBalances {
		yearMonth := balancesByDate.Date.Format("2006-01")
		if dataPointSet[yearMonth] == nil {
			dataPointSet[yearMonth] = DataPoint{}
		}
		for _, balance := range balancesByDate.Balances {
			dataPointSet[yearMonth]["liabilities"] = dataPointSet[yearMonth]["liabilities"] - balance.Balance
		}
	}

	for date := range dataPointSet {
		dataPointSet[date]["netWorth"] = dataPointSet[date]["assets"] + dataPointSet[date]["liabilities"]
	}

	tmplVariables.DataPoints = dataPointSet
	tmpl, err := template.New("netWorthDashboard").Parse(netWorthTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, tmplVariables)
	if err != nil {
		panic(err)
	}
}
