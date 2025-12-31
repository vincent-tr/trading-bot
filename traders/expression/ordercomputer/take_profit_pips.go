package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func TakeProfitPips(pips float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			entryPrice := ctx.EntryPrice()
			pipDistance := pips * pipSize

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.TakeProfit = entryPrice + pipDistance
				return nil

			case brokers.PositionDirectionShort:
				order.TakeProfit = entryPrice - pipDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "TakeProfitPips", formatter.FloatValue(pips))
		},
	)
}
