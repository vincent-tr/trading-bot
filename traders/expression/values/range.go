package values

import (
	"math"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

// RangeSize computes the range (high - low) over a specified number of historical candles.
// It returns the difference between the highest high and lowest low within the period.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeSize(20, 1) calculates the range from candles 1-20 bars ago (skipping the most recent completed candle).
func RangeSize(period int, offset int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			high := getRangeHigh(ctx, period, offset)
			low := getRangeLow(ctx, period, offset)
			return high - low
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeSize",
				formatter.IntValue(period),
				formatter.IntValue(offset),
			)
		},
	)
}

// RangeHigh returns the highest high price over a specified number of historical candles.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeHigh(20, 1) finds the highest price from candles 1-20 bars ago.
func RangeHigh(period int, offset int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			return getRangeHigh(ctx, period, offset)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeHigh",
				formatter.IntValue(period),
				formatter.IntValue(offset),
			)
		},
	)
}

// RangeLow returns the lowest low price over a specified number of historical candles.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeLow(20, 1) finds the lowest price from candles 1-20 bars ago.
func RangeLow(period int, offset int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			return getRangeLow(ctx, period, offset)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeLow",
				formatter.IntValue(period),
				formatter.IntValue(offset),
			)
		},
	)
}

func getRangeHigh(ctx context.TraderContext, period int, offset int) float64 {
	highs := ctx.HistoricalData().GetHighPrices()

	hightest := math.Inf(-1)

	for i := offset; i < period+offset; i++ {
		price := highs.At(i)
		if price > hightest {
			hightest = price
		}
	}

	return hightest
}

func getRangeLow(ctx context.TraderContext, period int, offset int) float64 {
	lows := ctx.HistoricalData().GetLowPrices()

	lowest := math.Inf(1)

	for i := offset; i < period+offset; i++ {
		price := lows.At(i)
		if price < lowest {
			lowest = price
		}
	}

	return lowest
}
