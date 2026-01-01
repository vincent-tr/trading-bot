package conditions

import (
	"math"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func BreakRangeHigh(loopbackPeriod int) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			high := ctx.HistoricalData().GetHighPrices()
			if high.Len() < loopbackPeriod+1 {
				return false
			}

			hightest := math.Inf(-1)
			for i := 1; i <= loopbackPeriod; i++ {
				price := high.Get(i)
				if price > hightest {
					hightest = price
				}
			}

			return ctx.HistoricalData().GetPrice() > hightest
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"BreakRangeHigh",
				formatter.IntValue(loopbackPeriod),
			)
		},
	)
}

func BreakRangeLow(loopbackPeriod int) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			low := ctx.HistoricalData().GetLowPrices()
			if low.Len() < loopbackPeriod+1 {
				return false
			}

			lowest := math.Inf(1)
			for i := 1; i <= loopbackPeriod; i++ {
				price := low.Get(i)
				if price < lowest {
					lowest = price
				}
			}

			return ctx.HistoricalData().GetPrice() < lowest
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"BreakRangeLow",
				formatter.IntValue(loopbackPeriod),
			)
		},
	)
}
