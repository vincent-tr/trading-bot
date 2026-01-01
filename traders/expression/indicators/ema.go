package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"

	"github.com/markcheno/go-talib"
)

func EMA(period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			closePrices := ctx.HistoricalData().GetClosePrices().All()
			return talib.Ema(closePrices, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "EMA", formatter.IntValue(period))
		},
	)
}
