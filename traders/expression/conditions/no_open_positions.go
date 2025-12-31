package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func NoOpenPositions() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return len(ctx.OpenPositions()) == 0
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "NoOpenPositions")
		},
	)
}
