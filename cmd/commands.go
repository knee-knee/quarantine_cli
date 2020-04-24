package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/urfave/cli"
)

func New(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:   "daysSinceStart",
			Usage:  "Tell you how long since this hell has started.",
			Action: DaysSinceStart,
		},
		{
			Name:   "checkUnemployment",
			Usage:  "Look at the current unemployment numbers",
			Action: CheckUnemployment,
		},
	}
}

func DaysSinceStart(c *cli.Context) {
	startDate := time.Date(2020, time.Month(3), 16, 0, 0, 0, 0, time.UTC)
	daysSince := (time.Now().Sub(startDate).Hours()) / 24                                // Dividing by 24 here to turn it into days.
	fmt.Printf("It has been %d days since the quarantine has started", int32(daysSince)) // Converting to an int here because I dont want decimals.
}

func CheckUnemployment(c *cli.Context) {
	res, err := http.Get("https://api.bls.gov/publicAPI/v1/timeseries/data/CEU0800000003")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var unemploymentRes UnemploymentRes
	if err := json.Unmarshal(body, &unemploymentRes); err != nil {
		log.Fatal(err)
	}

	var latestUnemploymentRate string
	for _, period := range unemploymentRes.Results.Series[0].Data { // We will only ever deal with one item and that is why we can get the first of the series.
		if period.Latest != nil && *period.Latest == "true" {
			latestUnemploymentRate = period.Value
		}
	}

	if latestUnemploymentRate == "" {
		log.Fatal("was not able to find the latest unemployment numbers")
	}

	fmt.Printf("Currently %s percent of the population is unemployed.\n", latestUnemploymentRate)
}

type UnemploymentRes struct {
	Status       string `json:"status"`
	ResponseTime int    `json:"responseTime"`
	Results      struct {
		Series []struct {
			SeriesID string `json:"seriesID"`
			Data     []struct {
				Year       string  `json:"year"`
				Period     string  `json:"period"`
				PeriodName string  `json:"periodName"`
				Latest     *string `json:"latest,omitempty"`
				Value      string  `json:"value"`
				Footnotes  []struct {
					Code string `json:"code"`
					Text string `json:"text"`
				} `json:"footnotes"`
			} `json:"data"`
		} `json:"series"`
	} `json:"Results"`
}
