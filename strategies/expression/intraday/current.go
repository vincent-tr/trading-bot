package intraday

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
		expression.Timeframe(brokers.Timeframe5Minutes),

		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),
				),
			),

			// LONG
			expression.LongTrigger(
				conditions.And(
					conditions.CrossAbove(
						indicators.EMA(20),
						indicators.EMA(50),
					),
					conditions.ValueBelow(
						indicators.RSI(14),
						values.Constant(70),
					),
				),
			),

			// SHORT
			expression.ShortTrigger(
				conditions.And(
					conditions.CrossBelow(
						indicators.EMA(20),
						indicators.EMA(50),
					),
					conditions.ValueAbove(
						indicators.RSI(14),
						values.Constant(30),
					),
				),
			),
		),

		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossDistance(
					values.Factor(
						indicators.ATR(14),
						1.0,
					),
				),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(1),
			),
		),

		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
