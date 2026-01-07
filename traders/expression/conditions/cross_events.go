package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/indicators"
)

// CrossAbove returns a condition that checks if indicator1 has crossed above indicator2.
func CrossAbove(indicator1, indicator2 indicators.Indicator) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			values1 := indicator1.Values(ctx)
			values2 := indicator2.Values(ctx)

			return values1.Previous() <= values2.Previous() && values1.Current() > values2.Current()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"CrossAbove",
				indicator1.Format(),
				indicator2.Format(),
			)
		},
	)
}

// CrossBelow returns a condition that checks if indicator1 has crossed below indicator2.
func CrossBelow(indicator1, indicator2 indicators.Indicator) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			values1 := indicator1.Values(ctx)
			values2 := indicator2.Values(ctx)

			return values1.Previous() >= values2.Previous() && values1.Current() < values2.Current()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"CrossBelow",
				indicator1.Format(),
				indicator2.Format(),
			)
		},
	)
}
