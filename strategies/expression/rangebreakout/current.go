package rangebreakout

import (
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

func Current() *expression.Configuration {

	const RangeDuration int = 10
	const ConfirmationDuration int = 3

	return expression.Builder(
		expression.HistorySize(250),
		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),
					// conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
					// conditions.ExcludeUKHolidays(),
					// conditions.ExcludeUSHolidays(),
					// conditions.Session(conditions.SessionLondon),
					// conditions.Session(conditions.SessionNewYork),

					// Volatility expansion
					conditions.ValueAbove(
						indicators.ATR(14),
						indicators.Mean(indicators.ATR(14), 10),
					),

					// Range must be tight
					conditions.ValueBelow(
						values.RangeSize(RangeDuration, values.Offset(ConfirmationDuration)),
						values.Factor(indicators.ATR(14), 1.2),
					),
				),
			),
			expression.LongTrigger(
				// Breakout acceptance
				conditions.ValueAbove(
					indicators.Min(indicators.Close(), ConfirmationDuration),
					values.RangeHigh(RangeDuration, values.Offset(ConfirmationDuration)),
				),
			),
			expression.ShortTrigger(
				// Breakout acceptance
				conditions.ValueBelow(
					indicators.Max(indicators.Close(), ConfirmationDuration),
					values.RangeLow(RangeDuration, values.Offset(ConfirmationDuration)),
				),
			),
		),
		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossFromRange(RangeDuration, 1, values.Offset(ConfirmationDuration)),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(1.5),
			),
		),
		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
