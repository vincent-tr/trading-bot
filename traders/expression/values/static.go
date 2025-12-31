package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Static(value float64) Value {
	return newValue(
		func(ctx context.TraderContext) float64 {
			return value
		},
		func() *formatter.FormatterNode {
			return formatter.FloatValue(value)
		},
	)
}
