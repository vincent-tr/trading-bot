package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/indicators"
)

// PriceAbove returns a condition that checks if the current price is above the given indicator value.
func PriceAbove(indicator indicators.Indicator) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			price := ctx.HistoricalData().GetPrice()
			indicatorValue := indicator.Get(ctx)

			return price > indicatorValue
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceAbove",
				indicator.Format(),
			)
		},
	)
}

// PriceBelow returns a condition that checks if the current price is below the given indicator value.
func PriceBelow(indicator indicators.Indicator) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			price := ctx.HistoricalData().GetPrice()
			indicatorValue := indicator.Get(ctx)

			return price < indicatorValue
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceBelow",
				indicator.Format(),
			)
		},
	)
}
