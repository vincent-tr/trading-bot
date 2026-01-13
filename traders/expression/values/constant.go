package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Constant(value float64) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			return value
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Constant", formatter.FloatValue(value))
		},
	)
}
