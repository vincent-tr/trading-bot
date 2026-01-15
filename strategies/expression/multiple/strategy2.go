package multiple

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

// STRATEGY 2 — Compression Breakout
func Strategy2() *expression.Configuration {
	return expression.Builder(
		expression.HistorySize(250),
		expression.Timeframe(brokers.Timeframe1Minute),

		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),
				),
			),

			// Break upper box
			expression.LongTrigger(
				conditions.PriceAbove(values.RangeHigh(50)),
			),

			// Break lower box
			expression.ShortTrigger(
				conditions.PriceBelow(values.RangeLow(50)),
			),
		),

		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossFromRange(20, 1.0),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(2.0),
			),
		),

		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
