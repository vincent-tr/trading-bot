package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"
	"trading-bot/brokers/backtesting"
	"trading-bot/common"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type bucket struct {
	date           string
	histCount      int
	dukasCount     int
	histMidPrice   float64
	dukasMidPrice  float64
	histSpread     float64
	dukasSpread    float64
	histTickCount  int
	dukasTickCount int
}

func main() {
	// Load both datasets for the same month
	month := common.NewMonth(2024, 1)

	// Calculate month boundaries for filtering
	monthStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // First day of next month

	histDataset, err := backtesting.LoadDataset(
		backtesting.HistData,
		month,
		month,
		"EURUSD",
	)
	if err != nil {
		panic(fmt.Errorf("failed to load histdata: %v", err))
	}

	dukasDataset, err := backtesting.LoadDataset(
		backtesting.Dukascopy,
		month,
		month,
		"EURUSD",
	)
	if err != nil {
		panic(fmt.Errorf("failed to load dukascopy: %v", err))
	}

	buckets := make([]*bucket, 0)
	bucketsMap := make(map[string]*bucket)

	// Fill buckets per hour only within the month boundaries
	for d := monthStart; d.Before(monthEnd); d = d.Add(time.Hour) {
		b := &bucket{date: d.Format("2006-01-02 15")}
		buckets = append(buckets, b)
		bucketsMap[b.date] = b
	}

	// Map histdata ticks to their respective buckets
	for tick := range histDataset.Ticks() {
		// Skip ticks outside the month boundaries
		timestamp := tick.GetTimestamp()
		if timestamp.Before(monthStart) || !timestamp.Before(monthEnd) {
			continue
		}

		tickDate := timestamp.Format("2006-01-02 15")
		bucket, ok := bucketsMap[tickDate]
		if !ok {
			continue
		}
		bucket.histCount++
		bucket.histTickCount++

		// Calculate mid price and spread
		midPrice := (tick.GetBid() + tick.GetAsk()) / 2
		spread := tick.GetAsk() - tick.GetBid()

		// Running average for the hour
		bucket.histMidPrice = (bucket.histMidPrice*float64(bucket.histTickCount-1) + midPrice) / float64(bucket.histTickCount)
		bucket.histSpread = (bucket.histSpread*float64(bucket.histTickCount-1) + spread) / float64(bucket.histTickCount)
	}

	// Map dukascopy ticks to their respective buckets
	for tick := range dukasDataset.Ticks() {
		// Skip ticks outside the month boundaries
		timestamp := tick.GetTimestamp()
		if timestamp.Before(monthStart) || !timestamp.Before(monthEnd) {
			continue
		}

		tickDate := timestamp.Format("2006-01-02 15")
		bucket, ok := bucketsMap[tickDate]
		if !ok {
			continue
		}
		bucket.dukasCount++
		bucket.dukasTickCount++

		// Calculate mid price and spread
		midPrice := (tick.GetBid() + tick.GetAsk()) / 2
		spread := tick.GetAsk() - tick.GetBid()

		// Running average for the hour
		bucket.dukasMidPrice = (bucket.dukasMidPrice*float64(bucket.dukasTickCount-1) + midPrice) / float64(bucket.dukasTickCount)
		bucket.dukasSpread = (bucket.dukasSpread*float64(bucket.dukasTickCount-1) + spread) / float64(bucket.dukasTickCount)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatal(err)
	}

	// Print some statistics to verify data differences
	fmt.Println("\nData Statistics:")
	totalHistTicks := 0
	totalDukasTicks := 0
	for _, b := range buckets {
		totalHistTicks += b.histCount
		totalDukasTicks += b.dukasCount
	}
	fmt.Printf("Total HistData ticks: %d\n", totalHistTicks)
	fmt.Printf("Total Dukascopy ticks: %d\n", totalDukasTicks)

	plotTickDensity(buckets)
	plotPriceAndSpread(buckets)
}

// weekendBand implements plot.Plotter
type weekendBand struct {
	xmin, xmax float64
	ymax       float64
}

func (w weekendBand) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	x0 := trX(w.xmin)
	x1 := trX(w.xmax)
	y0 := trY(0)
	y1 := trY(w.ymax)

	c.FillPolygon(color.RGBA{R: 200, G: 200, B: 200, A: 80},
		[]vg.Point{{X: x0, Y: y0}, {X: x1, Y: y0}, {X: x1, Y: y1}, {X: x0, Y: y1}})
}

func (w weekendBand) DataRange() (xmin, xmax, ymin, ymax float64) {
	return w.xmin, w.xmax, 0, w.ymax
}

func plotTickDensity(data []*bucket) {
	p := plot.New()
	p.Title.Text = "Tick Density Comparison (per Hour)"
	p.Y.Label.Text = "Tick Count"
	p.Legend.Top = true

	// Convert buckets → plot values
	histValues := make(plotter.Values, len(data))
	dukasValues := make(plotter.Values, len(data))
	labels := make([]string, len(data))
	maxY := 0.0

	for i, b := range data {
		histValues[i] = float64(b.histCount)
		dukasValues[i] = float64(b.dukasCount)

		if histValues[i] > maxY {
			maxY = histValues[i]
		}
		if dukasValues[i] > maxY {
			maxY = dukasValues[i]
		}

		t, _ := time.Parse("2006-01-02 15", b.date)
		if t.Hour() == 0 {
			prefix := ""
			if t.Weekday() == time.Monday {
				prefix = "| "
			}
			labels[i] = prefix + t.Format("01-02 (Mon)")
		}

		// Add shaded band for weekend hours
		if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			p.Add(weekendBand{
				xmin: float64(i) - 0.5,
				xmax: float64(i) + 0.5,
				ymax: maxY * 1.1,
			})
		}
	}

	// Histdata bars
	histBars, err := plotter.NewBarChart(histValues, vg.Points(8))
	if err != nil {
		log.Fatal(err)
	}
	histBars.Color = color.RGBA{R: 60, G: 120, B: 200, A: 200}
	histBars.LineStyle.Width = 0
	histBars.Offset = vg.Points(-4)
	p.Add(histBars)
	p.Legend.Add("HistData", histBars)

	// Dukascopy bars
	dukasBars, err := plotter.NewBarChart(dukasValues, vg.Points(8))
	if err != nil {
		log.Fatal(err)
	}
	dukasBars.Color = color.RGBA{R: 200, G: 60, B: 60, A: 200}
	dukasBars.LineStyle.Width = 0
	dukasBars.Offset = vg.Points(4)
	p.Add(dukasBars)
	p.Legend.Add("Dukascopy", dukasBars)

	p.NominalX(labels...)
	p.Add(plotter.NewGrid())

	// High-resolution canvas
	width := 2000
	height := 800
	canvas := vgimg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	p.Draw(draw.New(canvas))

	// Save as PNG
	f, err := os.Create("output/tick-density.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	png := vgimg.PngCanvas{Canvas: canvas}
	_, err = png.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tick density chart saved to output/tick-density.png")
}

func plotPriceAndSpread(data []*bucket) {
	pPrice := plot.New()
	pPrice.Title.Text = "Mid Price Comparison (per Hour)"
	pPrice.Y.Label.Text = "EUR/USD Price"
	pPrice.Legend.Top = true

	pSpread := plot.New()
	pSpread.Title.Text = "Spread Comparison (per Hour)"
	pSpread.Y.Label.Text = "Spread (pips x10000)"
	pSpread.Legend.Top = true

	// Convert buckets → plot points
	histPricePoints := make(plotter.XYs, 0)
	dukasPricePoints := make(plotter.XYs, 0)
	histSpreadPoints := make(plotter.XYs, 0)
	dukasSpreadPoints := make(plotter.XYs, 0)
	labels := make([]string, len(data))

	for i, b := range data {
		if b.histTickCount > 0 {
			histPricePoints = append(histPricePoints, plotter.XY{X: float64(i), Y: b.histMidPrice})
			histSpreadPoints = append(histSpreadPoints, plotter.XY{X: float64(i), Y: b.histSpread * 10000})
		}
		if b.dukasTickCount > 0 {
			dukasPricePoints = append(dukasPricePoints, plotter.XY{X: float64(i), Y: b.dukasMidPrice})
			dukasSpreadPoints = append(dukasSpreadPoints, plotter.XY{X: float64(i), Y: b.dukasSpread * 10000})
		}

		t, _ := time.Parse("2006-01-02 15", b.date)
		if t.Hour() == 0 {
			prefix := ""
			if t.Weekday() == time.Monday {
				prefix = "| "
			}
			labels[i] = prefix + t.Format("01-02 (Mon)")
		}
	}

	// Price chart - HistData line
	histPriceLine, err := plotter.NewLine(histPricePoints)
	if err != nil {
		log.Fatal(err)
	}
	histPriceLine.Color = color.RGBA{R: 60, G: 120, B: 200, A: 200} // Added transparency
	histPriceLine.Width = vg.Points(2)
	pPrice.Add(histPriceLine)
	pPrice.Legend.Add("HistData", histPriceLine)

	// Price chart - Dukascopy line
	dukasPriceLine, err := plotter.NewLine(dukasPricePoints)
	if err != nil {
		log.Fatal(err)
	}
	dukasPriceLine.Color = color.RGBA{R: 200, G: 60, B: 60, A: 200} // Added transparency
	dukasPriceLine.Width = vg.Points(2)
	pPrice.Add(dukasPriceLine)
	pPrice.Legend.Add("Dukascopy", dukasPriceLine)

	// Spread chart - HistData line
	histSpreadLine, err := plotter.NewLine(histSpreadPoints)
	if err != nil {
		log.Fatal(err)
	}
	histSpreadLine.Color = color.RGBA{R: 60, G: 120, B: 200, A: 200} // Added transparency
	histSpreadLine.Width = vg.Points(2)
	pSpread.Add(histSpreadLine)
	pSpread.Legend.Add("HistData", histSpreadLine)

	// Spread chart - Dukascopy line
	dukasSpreadLine, err := plotter.NewLine(dukasSpreadPoints)
	if err != nil {
		log.Fatal(err)
	}
	dukasSpreadLine.Color = color.RGBA{R: 200, G: 60, B: 60, A: 200} // Added transparency
	dukasSpreadLine.Width = vg.Points(2)
	pSpread.Add(dukasSpreadLine)
	pSpread.Legend.Add("Dukascopy", dukasSpreadLine)

	pPrice.NominalX(labels...)
	pPrice.Add(plotter.NewGrid())

	pSpread.NominalX(labels...)
	pSpread.Add(plotter.NewGrid())

	// High-resolution canvas
	width := 2000
	height := 800

	// Save price chart
	canvasPrice := vgimg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	pPrice.Draw(draw.New(canvasPrice))

	fPrice, err := os.Create("output/price-comparison.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fPrice.Close()

	pngPrice := vgimg.PngCanvas{Canvas: canvasPrice}
	_, err = pngPrice.WriteTo(fPrice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Price comparison chart saved to output/price-comparison.png")

	// Save spread chart
	canvasSpread := vgimg.New(vg.Points(float64(width)), vg.Points(float64(height)))
	pSpread.Draw(draw.New(canvasSpread))

	fSpread, err := os.Create("output/spread-comparison.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fSpread.Close()

	pngSpread := vgimg.PngCanvas{Canvas: canvasSpread}
	_, err = pngSpread.WriteTo(fSpread)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Spread comparison chart saved to output/spread-comparison.png")
}
