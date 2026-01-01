package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

// PriceAbove returns a condition that checks if the current price is above the given value.
func PriceAbove(value values.Value) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			price := ctx.HistoricalData().GetPrice()
			indicatorValue := value.Get(ctx)

			return price > indicatorValue
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceAbove",
				value.Format(),
			)
		},
	)
}

// PriceBelow returns a condition that checks if the current price is below the given value.
func PriceBelow(value values.Value) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			price := ctx.HistoricalData().GetPrice()
			valueValue := value.Get(ctx)

			return price < valueValue
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceBelow",
				value.Format(),
			)
		},
	)
}

// ValueAbove returns a condition that checks if valueA is above valueB.
func ValueAbove(valueA, valueB values.Value) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			valueA := valueA.Get(ctx)
			valueB := valueB.Get(ctx)

			return valueA > valueB
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"ValueAbove",
				valueA.Format(),
				valueB.Format(),
			)
		},
	)
}

// ValueBelow returns a condition that checks if valueA is below valueB.
func ValueBelow(valueA, valueB values.Value) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			valueA := valueA.Get(ctx)
			valueB := valueB.Get(ctx)

			return valueA < valueB
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"ValueBelow",
				valueA.Format(),
				valueB.Format(),
			)
		},
	)
}
