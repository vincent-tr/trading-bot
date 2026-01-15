package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

// N = loopbackPeriod
//
// NetMove   = | Close(0) − Close(N) |
//
// PathMove  = Sum( | Close(i) − Close(i+1) | for i=0..N-1 )
//
// Efficiency = NetMove / PathMove
func Efficiency(loopbackPeriod int) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			history := ctx.HistoricalData()
			closes := history.GetClosePrices()

			netMove := distance(closes.At(0), closes.At(loopbackPeriod))
			pathMove := 0.0
			for i := 0; i < loopbackPeriod; i++ {
				pathMove += distance(closes.At(i), closes.At(i+1))
			}

			return netMove / pathMove
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Efficiency", formatter.IntValue(loopbackPeriod))
		},
	)
}

func distance(a, b float64) float64 {
	val := a - b
	if val < 0 {
		return -val
	} else {
		return val
	}
}
