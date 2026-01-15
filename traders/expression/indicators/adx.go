package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"

	"github.com/markcheno/go-talib"
)

func ADX(period int) Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			history := ctx.HistoricalData()
			return talib.Adx(history.GetHighPrices().All(), history.GetLowPrices().All(), history.GetClosePrices().All(), period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "ADX", formatter.IntValue(period))
		},
	)
}
