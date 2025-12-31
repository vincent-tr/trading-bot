package conditions

import (
	"encoding/json"
	"fmt"
	"trading-bot/traders/modular/context"
	"trading-bot/traders/modular/formatter"
	"trading-bot/traders/modular/indicators"
	"trading-bot/traders/modular/marshal"
)

type Direction int

const (
	Above Direction = iota
	Below
)

func (d Direction) String() string {
	switch d {
	case Above:
		return "Above"
	case Below:
		return "Below"
	default:
		return "Unknown"
	}
}

// Threshold checks if the value of an indicator is above or below a specified threshold.
func Threshold(indicator indicators.Indicator, threshold float64, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			value := values[len(values)-1]

			switch direction {
			case Above:
				return value >= threshold
			case Below:
				return value <= threshold
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Threshold",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Value: %.2f", threshold)),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
		func() (string, any) {
			var directionStr string
			switch direction {
			case Above:
				directionStr = "above"
			case Below:
				directionStr = "below"
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}

			return "threshold", map[string]any{
				"indicator": marshal.ToJSON(indicator),
				"threshold": threshold,
				"direction": directionStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("threshold", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Threshold float64         `json:"threshold"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse Threshold parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		var direction Direction
		switch params.Direction {
		case "above":
			direction = Above
		case "below":
			direction = Below
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return Threshold(indicator, params.Threshold, direction), nil
	})
}

// PriceThreshold checks if the current price is above or below the value of an indicator.
func PriceThreshold(indicator indicators.Indicator, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) == 0 {
				return false
			}
			value := values[len(values)-1]
			entryPrice := ctx.EntryPrice()

			switch direction {
			case Above:
				return entryPrice >= value
			case Below:
				return entryPrice <= value
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("PriceThreshold",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
		func() (string, any) {
			var directionStr string
			switch direction {
			case Above:
				directionStr = "above"
			case Below:
				directionStr = "below"
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}

			return "priceThreshold", map[string]any{
				"indicator": marshal.ToJSON(indicator),
				"direction": directionStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("priceThreshold", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse PriceThreshold parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		var direction Direction
		switch params.Direction {
		case "above":
			direction = Above
		case "below":
			direction = Below
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return PriceThreshold(indicator, direction), nil
	})
}

type SlopeDirection int

const (
	SlopeRising SlopeDirection = iota
	SlopeFalling
)

func (d SlopeDirection) String() string {
	switch d {
	case SlopeRising:
		return "rising"
	case SlopeFalling:
		return "falling"
	default:
		return "unknown"
	}
}

// / Slope checks if the slope of an indicator is rising or falling.
func Slope(indicator indicators.Indicator, period int, direction SlopeDirection) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			values := indicator.Values(ctx)
			if len(values) < period+1 {
				return false
			}

			prevValue := values[len(values)-period-1]
			currValue := values[len(values)-1]

			switch direction {
			case SlopeRising:
				return currValue > prevValue
			case SlopeFalling:
				return currValue < prevValue
			default:
				panic(fmt.Sprintf("unknown slope direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Slope",
				indicator.Format(),
				formatter.Format(fmt.Sprintf("Period: %d", period)),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
		func() (string, any) {
			var directionStr string
			switch direction {
			case SlopeRising:
				directionStr = "rising"
			case SlopeFalling:
				directionStr = "falling"
			default:
				panic(fmt.Sprintf("unknown slope direction: %d", direction))
			}

			return "slope", map[string]any{
				"indicator": marshal.ToJSON(indicator),
				"period":    period,
				"direction": directionStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("slope", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			Indicator json.RawMessage `json:"indicator"`
			Period    int             `json:"period"`
			Direction string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse Slope parameters: %w", err)
		}

		indicator, err := indicators.FromJSON(params.Indicator)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicator: %w", err)
		}

		var direction SlopeDirection
		switch params.Direction {
		case "rising":
			direction = SlopeRising
		case "falling":
			direction = SlopeFalling
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return Slope(indicator, params.Period, direction), nil
	})
}

func Compare(indicatorA, indicatorB indicators.Indicator, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			valuesA := indicatorA.Values(ctx)
			valuesB := indicatorB.Values(ctx)

			if len(valuesA) == 0 || len(valuesB) == 0 {
				return false
			}

			valueA := valuesA[len(valuesA)-1]
			valueB := valuesB[len(valuesB)-1]

			switch direction {
			case Above:
				return valueA >= valueB
			case Below:
				return valueA <= valueB
			default:
				panic(fmt.Sprintf("unknown comparison direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Compare",
				indicatorA.Format(),
				indicatorB.Format(),
				formatter.Format(fmt.Sprintf("Direction: %s", direction.String())),
			)
		},
		func() (string, any) {
			var directionStr string
			switch direction {
			case Above:
				directionStr = "above"
			case Below:
				directionStr = "below"
			default:
				panic(fmt.Sprintf("unknown comparison direction: %d", direction))
			}

			return "compare", map[string]any{
				"indicatorA": marshal.ToJSON(indicatorA),
				"indicatorB": marshal.ToJSON(indicatorB),
				"direction":  directionStr,
			}
		},
	)
}

func init() {
	jsonParsers.RegisterParser("compare", func(arg json.RawMessage) (Condition, error) {
		var params struct {
			IndicatorA json.RawMessage `json:"indicatorA"`
			IndicatorB json.RawMessage `json:"indicatorB"`
			Direction  string          `json:"direction"`
		}

		if err := json.Unmarshal(arg, &params); err != nil {
			return nil, fmt.Errorf("failed to parse Compare parameters: %w", err)
		}

		indicatorA, err := indicators.FromJSON(params.IndicatorA)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicatorA: %w", err)
		}

		indicatorB, err := indicators.FromJSON(params.IndicatorB)
		if err != nil {
			return nil, fmt.Errorf("failed to parse indicatorB: %w", err)
		}

		var direction Direction
		switch params.Direction {
		case "above":
			direction = Above
		case "below":
			direction = Below
		default:
			return nil, fmt.Errorf("unknown direction: %s", params.Direction)
		}

		return Compare(indicatorA, indicatorB, direction), nil
	})
}
