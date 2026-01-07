package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/indicators"
)

func SlopeUp(indicator indicators.Indicator, period int) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if values.Len() < period+1 {
				return false
			}

			return values.Current() > values.At(period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"SlopeUp",
				indicator.Format(),
				formatter.IntValue(period),
			)
		},
	)
}

func SlopeDown(indicator indicators.Indicator, period int) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if values.Len() < period+1 {
				return false
			}

			return values.Current() < values.At(period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"SlopeDown",
				indicator.Format(),
				formatter.IntValue(period),
			)
		},
	)
}
