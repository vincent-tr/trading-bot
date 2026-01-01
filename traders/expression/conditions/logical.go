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

func Not(condition Condition) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return !condition.Execute(ctx)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Not",
				condition.Format(),
			)
		},
	)
}
