package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Or(conditions ...Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, condition := range conditions {
				if condition.Execute(ctx) {
					return true
				}
			}
			return false
		},
		func() *formatter.FormatterNode {
			return formatter.FunctionWithChildren(Package, "Or", conditions...)
		},
	)
}
