package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

func NormalizedSlope(indicator Indicator, period int) values.Value {
	return values.NewValue(
		func(ctx context.TraderContext) float64 {
			vals := indicator.Values(ctx)
			if vals.Len() < period+1 {
				panic("not enough data")
			}

			cur := vals.Current()
			prev := vals.At(period)
			atr := ATR(14).Values(ctx).Current()

			return (cur - prev) / atr
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"NormalizedSlope",
				indicator.Format(),
				formatter.IntValue(period),
			)
		},
	)
}
