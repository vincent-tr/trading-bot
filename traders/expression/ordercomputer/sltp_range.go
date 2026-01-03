package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

// StopLossFromRange sets the stop loss just outside the range high/low with a buffer
func StopLossFromRange(rangeLookback int, offset int, pipBuffer float64) OrderComputer {
	rangeLow := values.RangeLow(rangeLookback, offset)
	rangeHigh := values.RangeHigh(rangeLookback, offset)

	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			pipDistance := pipBuffer * pipSize

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.StopLoss = rangeLow.Get(ctx) - pipDistance
				return nil

			case brokers.PositionDirectionShort:
				order.StopLoss = rangeHigh.Get(ctx) + pipDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "StopLossFromRange",
				formatter.IntValue(rangeLookback),
				formatter.FloatValue(pipBuffer),
			)
		},
	)
}

// TakeProfitFromRange sets the take profit to price +/- size of range
func TakeProfitFromRange(rangeLookback int, offset int) OrderComputer {
	rangeSize := values.RangeSize(rangeLookback, offset)

	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			entryPrice := ctx.EntryPrice()
			rangeDistance := rangeSize.Get(ctx)

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.TakeProfit = entryPrice + rangeDistance
				return nil

			case brokers.PositionDirectionShort:
				order.TakeProfit = entryPrice - rangeDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "TakeProfitFromRange",
				formatter.IntValue(rangeLookback),
			)
		},
	)
}
