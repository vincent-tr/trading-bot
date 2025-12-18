package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
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

	// Print comprehensive statistics about data differences
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("DATA QUALITY ANALYSIS - HISTDATA vs DUKASCOPY")
	fmt.Println(strings.Repeat("=", 80))

	// Calculate totals and statistics
	totalHistTicks := 0
	totalDukasTicks := 0
	histHoursWithData := 0
	dukasHoursWithData := 0
	bothHoursWithData := 0
	histOnlyHours := 0
	dukasOnlyHours := 0

	var histPrices, dukasPrices, histSpreads, dukasSpreads []float64

	for _, b := range buckets {
		totalHistTicks += b.histCount
		totalDukasTicks += b.dukasCount

		hasHist := b.histCount > 0
		hasDukas := b.dukasCount > 0

		if hasHist {
			histHoursWithData++
			histPrices = append(histPrices, b.histMidPrice)
			histSpreads = append(histSpreads, b.histSpread)
		}
		if hasDukas {
			dukasHoursWithData++
			dukasPrices = append(dukasPrices, b.dukasMidPrice)
			dukasSpreads = append(dukasSpreads, b.dukasSpread)
		}
		if hasHist && hasDukas {
			bothHoursWithData++
		}
		if hasHist && !hasDukas {
			histOnlyHours++
		}
		if !hasHist && hasDukas {
			dukasOnlyHours++
		}
	}

	// Overall statistics
	fmt.Println("\n📊 TICK COUNT COMPARISON")
	fmt.Printf("  HistData:   %10d ticks (%6.2f%% of total)\n", totalHistTicks,
		float64(totalHistTicks)*100/float64(totalHistTicks+totalDukasTicks))
	fmt.Printf("  Dukascopy:  %10d ticks (%6.2f%% of total)\n", totalDukasTicks,
		float64(totalDukasTicks)*100/float64(totalHistTicks+totalDukasTicks))
	fmt.Printf("  Difference: %10d ticks (%.2fx more in %s)\n",
		abs(totalHistTicks-totalDukasTicks),
		float64(max(totalHistTicks, totalDukasTicks))/float64(min(totalHistTicks, totalDukasTicks)),
		ternary(totalHistTicks > totalDukasTicks, "HistData", "Dukascopy"))

	// Coverage statistics
	fmt.Println("\n⏰ TIME COVERAGE")
	fmt.Printf("  Total hours analyzed:     %d\n", len(buckets))
	fmt.Printf("  HistData hours:           %d (%5.2f%%)\n", histHoursWithData,
		float64(histHoursWithData)*100/float64(len(buckets)))
	fmt.Printf("  Dukascopy hours:          %d (%5.2f%%)\n", dukasHoursWithData,
		float64(dukasHoursWithData)*100/float64(len(buckets)))
	fmt.Printf("  Both sources:             %d (%5.2f%%)\n", bothHoursWithData,
		float64(bothHoursWithData)*100/float64(len(buckets)))
	fmt.Printf("  HistData only:            %d\n", histOnlyHours)
	fmt.Printf("  Dukascopy only:           %d\n", dukasOnlyHours)

	// Density statistics
	if histHoursWithData > 0 && dukasHoursWithData > 0 {
		avgHistTicksPerHour := float64(totalHistTicks) / float64(histHoursWithData)
		avgDukasTicksPerHour := float64(totalDukasTicks) / float64(dukasHoursWithData)

		fmt.Println("\n📈 AVERAGE TICK DENSITY (per hour with data)")
		fmt.Printf("  HistData:   %8.2f ticks/hour\n", avgHistTicksPerHour)
		fmt.Printf("  Dukascopy:  %8.2f ticks/hour\n", avgDukasTicksPerHour)
		fmt.Printf("  Difference: %8.2f ticks/hour (%.2f%%)\n",
			avgDukasTicksPerHour-avgHistTicksPerHour,
			(avgDukasTicksPerHour-avgHistTicksPerHour)*100/avgHistTicksPerHour)
	}

	// Price statistics
	if len(histPrices) > 0 && len(dukasPrices) > 0 {
		histMinPrice, histMaxPrice := minMax(histPrices)
		dukasMinPrice, dukasMaxPrice := minMax(dukasPrices)
		histAvgSpread := average(histSpreads)
		dukasAvgSpread := average(dukasSpreads)

		fmt.Println("\n💰 PRICE RANGE")
		fmt.Printf("  HistData:   %.5f - %.5f (range: %.5f)\n", histMinPrice, histMaxPrice, histMaxPrice-histMinPrice)
		fmt.Printf("  Dukascopy:  %.5f - %.5f (range: %.5f)\n", dukasMinPrice, dukasMaxPrice, dukasMaxPrice-dukasMinPrice)

		fmt.Println("\n📏 AVERAGE SPREAD")
		fmt.Printf("  HistData:   %.5f (%.2f pips)\n", histAvgSpread, histAvgSpread*10000)
		fmt.Printf("  Dukascopy:  %.5f (%.2f pips)\n", dukasAvgSpread, dukasAvgSpread*10000)
		fmt.Printf("  Difference: %.5f (%.2f pips)\n", dukasAvgSpread-histAvgSpread, (dukasAvgSpread-histAvgSpread)*10000)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📁 Generating charts...")
	fmt.Println(strings.Repeat("=", 80))

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

// Helper functions for statistics
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func minMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
