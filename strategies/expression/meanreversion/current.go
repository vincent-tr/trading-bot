package meanreversion

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

func Current() *expression.Configuration {
	return expression.Builder(
		expression.HistorySize(250),
		expression.Timeframe(brokers.Timeframe1Minute),

		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),

					// Market must be non-trending:
					// Range must be tight relative to volatility
					conditions.ValueBelow(
						values.RangeSize(20),
						values.Factor(indicators.ATR(14), 1.20),
					),
				),
			),

			// LONG
			expression.LongTrigger(
				conditions.ValueBelow(
					indicators.RSI(14),
					values.Constant(35),
				),
			),

			// SHORT
			expression.ShortTrigger(
				conditions.ValueAbove(
					indicators.RSI(14),
					values.Constant(65),
				),
			),
		),

		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossDistance(
					values.Factor(
						indicators.ATR(10),
						1.0,
					),
				),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(0.8),
			),
		),

		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
