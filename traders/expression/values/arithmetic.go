package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Add(valueA, valueB Value) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			return valueA.Get(ctx) + valueB.Get(ctx)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Add", valueA.Format(), valueB.Format())
		},
	)
}

func Subtract(valueA, valueB Value) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			return valueA.Get(ctx) - valueB.Get(ctx)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Subtract", valueA.Format(), valueB.Format())
		},
	)
}
