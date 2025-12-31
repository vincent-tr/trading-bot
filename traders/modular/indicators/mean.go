package indicators

import (
	"fmt"
	"trading-bot/traders/modular/context"
	"trading-bot/traders/modular/formatter"
	"trading-bot/traders/modular/marshal"
)

// Mean computes the mean of the given indicator over the specified period.
func Mean(indicator Indicator, period int) Indicator {
	return newIndicator(
		func(ctx context.TraderContext) []float64 {
			values := indicator.Values(ctx)
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
			return formatter.Format("Mean",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Period: %d", period)),
			)
		},
		func() (string, any) {
			return "mean", map[string]any{
				"indicator": marshal.ToJSON(indicator),
				"period":    period,
			}
		},
	)
}
