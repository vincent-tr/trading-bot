package conditions

import (
	"fmt"
	"trading-bot/brokers"
	"trading-bot/traders/modular/context"
	"trading-bot/traders/modular/formatter"
)

type CandleProperty int

const (
	OpenPrice CandleProperty = iota
	ClosePrice
	HighPrice
	LowPrice
)

func Candle(referenceIndex int, referenceProperty CandleProperty, testIndex int, testProperty CandleProperty, direction Direction) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			referenceCandle := ctx.HistoricalData().GetCandle(referenceIndex)
			testCandle := ctx.HistoricalData().GetCandle(testIndex)

			referenceValue := getCandleValue(&referenceCandle, referenceProperty)
			testValue := getCandleValue(&testCandle, testProperty)

			switch direction {
			case Above:
				return testValue > referenceValue
			case Below:
				return testValue < referenceValue
			default:
				panic(fmt.Sprintf("unknown threshold direction: %d", direction))
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Format("Candle",
				formatter.Format(fmt.Sprintf("ReferenceIndex: %d", referenceIndex)),
				formatter.Format(fmt.Sprintf("ReferenceProperty: %d", referenceProperty)),
				formatter.Format(fmt.Sprintf("TestIndex: %d", testIndex)),
				formatter.Format(fmt.Sprintf("TestProperty: %d", testProperty)),
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

			return "candle", map[string]any{
				"referenceIndex":    referenceIndex,
				"referenceProperty": referenceProperty,
				"testIndex":         testIndex,
				"testProperty":      testProperty,
				"direction":         directionStr,
			}
		},
	)
}

func getCandleValue(candle *brokers.Candle, property CandleProperty) float64 {
	switch property {
	case OpenPrice:
		return candle.Open
	case ClosePrice:
		return candle.Close
	case HighPrice:
		return candle.High
	case LowPrice:
		return candle.Low
	default:
		panic("unknown candle property")
	}
}
