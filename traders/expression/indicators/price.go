package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Close() Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			return ctx.HistoricalData().GetClosePrices().All()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Close")
		},
	)
}

func High() Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			return ctx.HistoricalData().GetHighPrices().All()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "High")
		},
	)
}

func Low() Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			return ctx.HistoricalData().GetLowPrices().All()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Low")
		},
	)
}
