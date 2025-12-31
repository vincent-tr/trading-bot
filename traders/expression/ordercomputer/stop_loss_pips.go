package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func StopLossPips(pips float64) OrderComputer {
	return newOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			entryPrice := ctx.EntryPrice()
			pipDistance := pips * pipSize

			switch order.Direction {
			case brokers.PositionDirectionLong:
				order.StopLoss = entryPrice - pipDistance
				return nil

			case brokers.PositionDirectionShort:
				order.StopLoss = entryPrice + pipDistance
				return nil

			default:
				panic("invalid position type")
			}
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "StopLossPips", formatter.FloatValue(pips))
		},
	)
}
