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

type netWorthDataPoint struct {
	FormattedYearMonth string
	AssetsTotal        float32
	LiabilitiesTotal   float32
}

func StartApiServer(ctx context.Context) {
	http.HandleFunc("/", dashboardHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dashboardHandler(w http.ResponseWriter, req *http.Request) {

	/*

		What are we doing? We need a list of year-months (like ['2024-05', '2024-06'] and a list of corresponding data points [{assets: 123, liabilities 10}, ...])


		1. Pull assets, liabilities and calculate net worth
		2. Build a templateDataStruct?
		3.
	*/

	// hardcoded for prototyping, need to get this from the UI eventually
	dateRange := []time.Time{}
	jan, _ := time.Parse("2006-01", "2024-01")
	feb, _ := time.Parse("2006-01", "2024-02")
	mar, _ := time.Parse("2006-01", "2024-03")
	apr, _ := time.Parse("2006-01", "2024-04")
	may, _ := time.Parse("2006-01", "2024-05")
	jun, _ := time.Parse("2006-01", "2024-06")
	dateRange = append(dateRange, jan)
	dateRange = append(dateRange, feb)
	dateRange = append(dateRange, mar)
	dateRange = append(dateRange, apr)
	dateRange = append(dateRange, may)
	dateRange = append(dateRange, jun)

	br := models.GetBalanceRepository()
	assetBalances := br.GetBalancesOfAllAssets(context.TODO(), "2024-01", "2024-06")
	liabilityBalances := br.GetBalancesOfAllLiabilities(context.TODO(), "2024-01", "2024-06")

	// totalAssetsByMonth := map[time.Time]float32{}
	// totalLiabilitiesByMonth := map[time.Time]float32{}

	dataPointSet := map[string]DataPoint{}
	// dataPointSet["2024-01"] = dataPoint{
	// 	"assets": 123.0,
	// }

	for _, date := range dateRange {
		// fmt.Println(yearMonth)
		// fmt.Println(totalAssetsByMonth)
		// totalAssetsByMonth[yearMonth] = 0.0
		// totalLiabilitiesByMonth[yearMonth] = 0.0
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
				dataPointSet[yearMonth]["liabilities"] = dataPointSet[yearMonth]["liabilities"] + balance.Balance
			}
		}
	}

	// fmt.Printf("asset totals: %v\n", totalAssetsByMonth)
	// fmt.Printf("liability totals: %v\n", totalAssetsByMonth)

	tmplVariables := dashboardTmplVariables{
		DataPoints: dataPointSet,
	}
	tmpl, err := template.New("test").Parse(dashboardTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, tmplVariables)
	if err != nil {
		panic(err)
	}
	// io.WriteString(w, dashboardTmpl)
}
