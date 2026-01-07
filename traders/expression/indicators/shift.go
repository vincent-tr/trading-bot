package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

// Shift returns a Value that shifts the given indicator by the specified offset in the past.
func Shift(indicator Indicator, offset int) values.Value {
	return values.NewValue(
		func(ctx context.TraderContext) float64 {
			return indicator.Values(ctx).At(offset)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Shift", indicator.Format(), formatter.IntValue(offset))
		},
	)
}
