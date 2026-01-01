package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Mean(indicator Indicator, period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			values := indicator.Values(ctx).All()
			if len(values) < period {
				return []float64{}
			}

			meanValues := make([]float64, len(values)-period+1)

			var sum float64
			// Calculate sum for first window
			for i := 0; i < period; i++ {
				sum += values[i]
			}
			meanValues[0] = sum / float64(period)

			// Slide the window
			for i := 1; i < len(meanValues); i++ {
				sum -= values[i-1]
				sum += values[i+period-1]
				meanValues[i] = sum / float64(period)
			}

			return meanValues
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Mean", indicator.Format(), formatter.IntValue(period))
		},
	)
}
