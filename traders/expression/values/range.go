package values

import (
	"math"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

// RangeSize computes the range (high - low) of the given period.
// Range values never include the current candle.
func RangeSize(period int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			high := getRangeHigh(ctx, period)
			low := getRangeLow(ctx, period)
			return high - low
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeSize",
				formatter.IntValue(period),
			)
		},
	)
}

func RangeHigh(period int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			return getRangeHigh(ctx, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeHigh",
				formatter.IntValue(period),
			)
		},
	)
}

func RangeLow(period int) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			return getRangeLow(ctx, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeLow",
				formatter.IntValue(period),
			)
		},
	)
}

func getRangeHigh(ctx context.TraderContext, period int) float64 {
	highs := ctx.HistoricalData().GetHighPrices()

	hightest := math.Inf(-1)

	for i := 1; i <= period; i++ {
		price := highs.At(i)
		if price > hightest {
			hightest = price
		}
	}

	return hightest
}

func getRangeLow(ctx context.TraderContext, period int) float64 {
	lows := ctx.HistoricalData().GetLowPrices()

	lowest := math.Inf(1)

	for i := 1; i <= period; i++ {
		price := lows.At(i)
		if price < lowest {
			lowest = price
		}
	}

	return lowest
}
