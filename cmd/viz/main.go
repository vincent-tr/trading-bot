package main

import (
	"fmt"
	"net/http"
	"time"
	"trading-bot/brokers"
	"trading-bot/brokers/backtesting"
	"trading-bot/common"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Candle struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
}

var candles []Candle

func main() {
	dataset, err := backtesting.LoadDataset(
		backtesting.HistData,
		common.NewMonth(2024, 1),
		common.NewMonth(2024, 12),
		"EURUSD",
	)

	if err != nil {
		panic(err)
	}

	brokerConfig := &backtesting.Config{
		// For backtesting, we assume a lot size of 1 for simplicity.
		// In a real broker, this would be the number of units per lot.
		// Not that using IG broker, EUR/USD Mini has also a size of 1.
		LotSize: 1,

		// Leverage is the ratio of the amount of capital that a trader must put up to open a position.
		// For example, if the leverage is 30, it means that for every 1 unit of capital,
		// the trader can control 30 units of the asset.
		// This is a common leverage ratio in forex trading.
		Leverage: 30.0,

		InitialCapital: 100000,
	}

	// Create a backtesting broker to generate M1 candles
	broker, err := backtesting.NewBroker(brokerConfig, dataset)
	if err != nil {
		panic(err)
	}

	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		candles = append(candles, Candle{
			Timestamp: broker.GetCurrentTime(),
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
		})
	})

	// Run the broker to generate all candles
	if err := broker.Run(); err != nil {
		panic(err)
	}

	fmt.Printf("Got %d candles\n", len(candles))

	http.HandleFunc("/", handleChart)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleChart(w http.ResponseWriter, r *http.Request) {
	start := 0
	size := 500

	if v := r.URL.Query().Get("start"); v != "" {
		fmt.Sscanf(v, "%d", &start)
	}
	if v := r.URL.Query().Get("size"); v != "" {
		fmt.Sscanf(v, "%d", &size)
	}

	if start < 0 {
		start = 0
	}
	if start >= len(candles) {
		start = len(candles) - size
	}
	if start < 0 {
		start = 0
	}
	end := start + size
	if end > len(candles) {
		end = len(candles)
	}
	window := candles[start:end]

	// Convert to echarts format
	klineData := make([]opts.KlineData, len(window))
	xAxis := make([]string, len(window))

	for i, candle := range window {
		klineData[i] = opts.KlineData{
			Value: [4]float64{
				candle.Open,
				candle.Close,
				candle.Low,
				candle.High,
			},
		}
		xAxis[i] = candle.Timestamp.Format("2006-01-02 15:04")
	}

	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "EUR/USD M1 Candles",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 0,
			End:   100,
		}),
	)

	kline.SetXAxis(xAxis).AddSeries("EUR/USD", klineData)

	prevStart := start - size
	nextStart := start + size

	if prevStart < 0 {
		prevStart = 0
	}
	if nextStart >= len(candles) {
		nextStart = start // disable forward
	}

	fmt.Fprintf(w, `
<html>
<head>
    <meta charset="utf-8">
    <title>EUR/USD Debug</title>
</head>
<body>

<div style="margin-bottom:10px;">
    <a href="/?start=%d&size=%d">⬅ Previous</a>
    |
    <a href="/?start=%d&size=%d">Next ➡</a>
    <span style="margin-left:20px;">
        Showing candles %d → %d / %d
    </span>
</div>
`,
		prevStart, size,
		nextStart, size,
		start, end, len(candles),
	)

	page := components.NewPage()
	page.AddCharts(kline)
	page.Render(w)

	fmt.Fprint(w, "</body></html>")
}
