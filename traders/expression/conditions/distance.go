package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/indicators"
)

// DistanceIncreasing returns a condition that checks if the distance between the price and the indicator is increasing.
func DistanceIncreasing(indicator indicators.Indicator) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			prices := ctx.HistoricalData().GetClosePrices()
			values := indicator.Values(ctx)
			if values.Len() < 2 {
				return false
			}

			currentDistance := abs(prices.Current() - values.Current())
			previousDistance := abs(prices.Previous() - values.Previous())

			return currentDistance > previousDistance
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"DistanceIncreasing",
				indicator.Format(),
			)
		},
	)
}

func DistanceDecreasing(indicator indicators.Indicator) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			prices := ctx.HistoricalData().GetClosePrices()
			values := indicator.Values(ctx)
			if values.Len() < 2 {
				return false
			}

			currentDistance := abs(prices.Current() - values.Current())
			previousDistance := abs(prices.Previous() - values.Previous())

			return currentDistance < previousDistance
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"DistanceDecreasing",
				indicator.Format(),
			)
		},
	)
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
