package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func HistoryUsable() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return ctx.HistoricalData().IsUsable()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "HistoryUsable")
		},
	)
}
