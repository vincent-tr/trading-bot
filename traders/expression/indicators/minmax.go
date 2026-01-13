package indicators

import (
	"math"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

func Min(indicator Indicator, period int) values.Value {
	return values.NewValue(
		func(ctx context.TraderContext) float64 {
			data := indicator.Values(ctx)

			if data.Len() < period {
				panic("not enough data for Min indicator")
			}

			min := math.Inf(1)

			for i := 0; i < period; i++ {
				curr := data.At(i)
				if curr < min {
					min = curr
				}
			}

			return min
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"Min",
				indicator.Format(),
				formatter.IntValue(period),
			)
		},
	)
}

func Max(indicator Indicator, period int) values.Value {
	return values.NewValue(
		func(ctx context.TraderContext) float64 {
			data := indicator.Values(ctx)

			if data.Len() < period {
				panic("not enough data for Max indicator")
			}

			max := math.Inf(-1)

			for i := 0; i < period; i++ {
				curr := data.At(i)
				if curr > max {
					max = curr
				}
			}

			return max
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"Max",
				indicator.Format(),
				formatter.IntValue(period),
			)
		},
	)
}
