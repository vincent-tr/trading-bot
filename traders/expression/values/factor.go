package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Factor(value Value, factor float64) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			return value.Get(ctx) * factor
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Factor", value.Format(), formatter.FloatValue(factor))
		},
	)
}
