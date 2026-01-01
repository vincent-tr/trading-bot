package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"

	"github.com/markcheno/go-talib"
)

func ATR(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			history := ctx.HistoricalData()
			return talib.Atr(history.GetHighPrices().All(), history.GetLowPrices().All(), history.GetClosePrices().All(), period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "ATR", formatter.IntValue(period))
		},
	)
}
