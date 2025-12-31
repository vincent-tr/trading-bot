package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func And(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if !condition.Execute(ctx) {
					return false
				}
			}
			return true
		},
		func() *formatter.FormatterNode {
			return formatter.FunctionWithChildren(Package, "And", conditions...)
		},
	)
}
