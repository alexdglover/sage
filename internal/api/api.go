package api

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/alexdglover/sage/internal/models"
)

//go:embed dashboard.html
var dashboardTmpl string

// A dataPoint will be either 'asset': <some float> or 'liablities': <some float>
type DataPoint map[string]float32

type dashboardTmplVariables struct {
	DataPoints map[string]DataPoint
}

func StartApiServer(ctx context.Context) {
	http.HandleFunc("/", dashboardHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: Consider moving this into a service class that returns just the data needed
func dashboardHandler(w http.ResponseWriter, req *http.Request) {
	// hardcoded for prototyping, need to get this from the UI eventually

	dateRange := []time.Time{}

	// Let's assume we'll always get the end date, as this will be today
	endDate, err := time.Parse("2006-01", "2024-07")
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}
	// And we want the last 6 months (this will be variable, but let's
	// assume it will be the last N months)
	startDate := endDate.AddDate(0, -5, 0)

	// Iterate over each month between start and end dates
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 1, 0) {
		dateRange = append(dateRange, date)
	}

	br := models.GetBalanceRepository()
	assetBalances := br.GetBalancesOfAllAssets(context.TODO(), startDate.Format("2006-01"), endDate.Format("2006-01"))
	liabilityBalances := br.GetBalancesOfAllLiabilities(context.TODO(), startDate.Format("2006-01"), endDate.Format("2006-01"))

	dataPointSet := map[string]DataPoint{}

	for _, date := range dateRange {
		yearMonth := date.Format("2006-01")
		dataPointSet[yearMonth] = DataPoint{}

		for _, balance := range assetBalances {
			esd, err := time.Parse("2006-01-02", balance.EffectiveStartDate)
			if err != nil {
				fmt.Println("Could not parse time:", err)
			}
			if esd.Year() == date.Year() && esd.Month() == date.Month() {
				dataPointSet[yearMonth]["assets"] = dataPointSet[yearMonth]["assets"] + balance.Balance
			}
		}

		for _, balance := range liabilityBalances {
			esd, err := time.Parse("2006-01-02", balance.EffectiveStartDate)
			if err != nil {
				fmt.Println("Could not parse time:", err)
			}
			if esd.Year() == date.Year() && esd.Month() == date.Month() {
				dataPointSet[yearMonth]["liabilities"] = dataPointSet[yearMonth]["liabilities"] - balance.Balance
			}
		}
	}

	for date := range dataPointSet {
		dataPointSet[date]["netWorth"] = dataPointSet[date]["assets"] + dataPointSet[date]["liabilities"]
	}

	tmplVariables := dashboardTmplVariables{
		DataPoints: dataPointSet,
	}
	tmpl, err := template.New("netWorthDashboard").Parse(dashboardTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, tmplVariables)
	if err != nil {
		panic(err)
	}
}
