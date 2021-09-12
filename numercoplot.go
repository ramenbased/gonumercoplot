package main

import (
	"context"
	"math"
	"os"
	"strconv"
	"time"
	"unicode"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/wcharczuk/go-chart"
)

var bids []float64
var asks []float64
var spreads []float64
var spot []float64
var dates []time.Time
var tweeterr bool

func main() {
	numercoScrape()
}

func numercoScrape() {
	for {
		tweeterr = false
		scraper := twitterscraper.New()

		bids = nil
		asks = nil
		spreads = nil
		spot = nil
		dates = nil

		for tweet := range scraper.GetTweets(context.Background(), "Numerco", 10000) {
			if tweet.Error != nil {
				tweeterr = true
			}

			if len(tweet.Text) > 23 {
				bidask := tweet.Text[14:18] + tweet.Text[19:23]
				check := 0
				for _, unirunes := range bidask {
					if unicode.IsNumber(unirunes) {
						check += 1
					}
					if check == 8 {
						bid, _ := strconv.ParseFloat(bidask[:2]+"."+bidask[2:4], 64)
						ask, _ := strconv.ParseFloat(bidask[4:6]+"."+bidask[6:8], 64)

						bids = append(bids, bid)
						asks = append(asks, ask)
						spreads = append(spreads, math.Round((ask-bid)*100)/100)
						spot = append(spot, math.Round(((bid+ask)/2)*100)/100)
						dates = append(dates, tweet.TimeParsed)
					}
				}
			}
		}
		ts1 := chart.TimeSeries{

			Style: chart.Style{
				StrokeColor: chart.GetDefaultColor(2),
			},
			Name:    "ASK",
			XValues: dates,
			YValues: asks,
		}

		ts2 := chart.TimeSeries{
			Style: chart.Style{
				StrokeColor: chart.GetDefaultColor(0),
			},
			Name:    "BID",
			XValues: dates,
			YValues: bids,
		}

		graph := chart.Chart{

			XAxis: chart.XAxis{
				Name:           "data: @numerco twitter",
				ValueFormatter: chart.TimeDateValueFormatter,
			},

			YAxis: chart.YAxis{
				Name: "USD per lbs",
			},

			Series: []chart.Series{
				ts1,
				ts2,
				chart.AnnotationSeries{
					Annotations: []chart.Value2{
						{
							XValue: chart.TimeToFloat64(dates[0]),
							YValue: asks[0],
							Label:  strconv.FormatFloat(asks[0], 'f', 2, 64),
						},
					},
				},
				chart.AnnotationSeries{
					Style: chart.Style{
						StrokeColor: chart.GetDefaultColor(0),
					},

					Annotations: []chart.Value2{
						{
							XValue: chart.TimeToFloat64(dates[0]),
							YValue: bids[0],
							Label:  strconv.FormatFloat(bids[0], 'f', 2, 64),
						},
					},
				},
			},
		}

		graph.Elements = []chart.Renderable{
			chart.Legend(&graph),
		}
		os.Remove("output1.png")
		f, _ := os.Create("output1.png")
		defer f.Close()
		graph.Render(chart.PNG, f)
		time.Sleep(time.Duration(1) * time.Hour)
	}
}
