package conditions

import (
	"encoding/json"
	"fmt"
	"trading-bot/traders/modular/context"
	"trading-bot/traders/modular/formatter"
	"trading-bot/traders/modular/indicators"
)

// check |close(now) - indicator(now)| >or< |close(now-period) - indicator(now-period)|
func DistancePriceToPrevious(indicator indicators.Indicator, period int, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			history := ctx.HistoricalData()
			values := indicator.Values(ctx)
			if len(values) < period+1 || !history.IsUsable() {
				return false
			}

			prices := history.GetClosePrices()

			currentDistance := abs(prices[len(prices)-1] - values[len(values)-1])
			previousDistance := abs(prices[len(prices)-period-1] - values[len(values)-period-1])

			switch direction {
			case Above:
				return currentDistance > previousDistance
			case Below:
				return currentDistance < previousDistance
			default:
				panic("unknown direction")
			}
		},
		func() *formatter.FormatterNode {
			var dirStr string
			switch direction {
			case Above:
				dirStr = ">"
			case Below:
				dirStr = "<"
			default:
				dirStr = "unknown"
			}

			return formatter.Format("DistancePriceToPrevious",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Period: %d", period)),
				formatter.Format(fmt.Sprintf("Direction: %s", dirStr)),
			)
		},
		func() (string, any) {
			return "distance_price_to_previous", map[string]any{
				"indicator": indicator,
				"period":    period,
				"direction": direction,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("distance_price_to_previous", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator indicators.Indicator `json:"indicator"`
			Period    int                  `json:"period"`
			Direction Direction            `json:"direction"`
		}
		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse DistancePriceToPrevious params: %w", err)
		}

		return DistancePriceToPrevious(params.Indicator, params.Period, params.Direction), nil
	})
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
