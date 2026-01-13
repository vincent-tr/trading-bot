package multiple

import (
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

// STRATEGY 1 — Compression Fade (RSI Mean Reversion)
func Strategy1() *expression.Configuration {
	return expression.Builder(
		expression.HistorySize(250),

		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),
					// Regime: Compression
					conditions.ValueBelow(
						values.RangeSize(50),
						values.Factor(indicators.ATR(50), 1.2),
					),
				),
			),

			// LONG: RSI washed out
			expression.LongTrigger(
				conditions.ValueBelow(indicators.RSI(14), values.Constant(25)),
			),

			// SHORT: RSI stretched
			expression.ShortTrigger(
				conditions.ValueAbove(indicators.RSI(14), values.Constant(75)),
			),
		),

		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossDistance(
					values.Factor(indicators.ATR(14), 0.8),
				),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(1.0),
			),
		),

		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
