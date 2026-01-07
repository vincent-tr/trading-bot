package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func BullishCandle() Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			candle := ctx.HistoricalData().GetCandle(0)
			return candle.Close > candle.Open
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"BullishCandle",
			)
		},
	)
}

func BearishCandle() Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			candle := ctx.HistoricalData().GetCandle(0)
			return candle.Close < candle.Open
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"BearishCandle",
			)
		},
	)
}
