package multiple

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

// STRATEGY 3 — Trend Pullback
func Strategy3() *expression.Configuration {
	return expression.Builder(
		expression.HistorySize(250),
		expression.Timeframe(brokers.Timeframe1Minute),

		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),

					// Regime: Trend
					conditions.ValueAbove(
						values.Abs(
							values.Subtract(indicators.Close(), indicators.SMA(50)),
						),
						values.Factor(indicators.ATR(50), 1.5),
					),
				),
			),

			// Buy pullback in uptrend
			expression.LongTrigger(
				conditions.And(
					conditions.ValueAbove(indicators.RSI(14), values.Constant(50)),
					conditions.ValueBelow(indicators.RSI(14), values.Constant(60)),
				),
			),

			// Sell pullback in downtrend
			expression.ShortTrigger(
				conditions.And(
					conditions.ValueBelow(indicators.RSI(14), values.Constant(50)),
					conditions.ValueAbove(indicators.RSI(14), values.Constant(40)),
				),
			),
		),

		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossDistance(
					values.Factor(indicators.ATR(14), 1.2),
				),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(2.5),
			),
		),

		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
