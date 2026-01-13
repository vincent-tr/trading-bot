package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

func StopLossDistance(value values.Value) OrderComputer {
	return NewOrderComputer(
		func(ctx context.TraderContext, order *brokers.Order) error {
			pipDistance := value.Get(ctx)
			entryPrice := ctx.EntryPrice()

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
			return formatter.Function(Package, "StopLossDistance", value.Format())
		},
	)
}
