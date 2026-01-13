package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"

	"github.com/markcheno/go-talib"
)

func RSI(period int) Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			closePrices := ctx.HistoricalData().GetClosePrices().All()
			return talib.Rsi(closePrices, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "RSI", formatter.IntValue(period))
		},
	)
}
