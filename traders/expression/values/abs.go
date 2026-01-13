package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Abs(value Value) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			v := value.Get(ctx)
			if v < 0 {
				return -v
			}
			return v
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Abs", value.Format())
		},
	)
}
