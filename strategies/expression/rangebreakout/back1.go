package rangebreakout

import (
	"trading-bot/traders/expression"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/expression/values"
)

func Back1() *expression.Configuration {

	const RangeDuration int = 10

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
					//conditions.ValueAbove(
					//	indicators.ATR(14),
					//	indicators.Mean(indicators.ATR(14), 10),
					//),

					// Range must be tight
					conditions.ValueBelow(
						values.RangeSize(RangeDuration),
						indicators.ATR(14),
					),
				),
			),
			expression.LongTrigger(
				conditions.PriceAbove(
					values.RangeHigh(RangeDuration),
				),
			),
			expression.ShortTrigger(
				conditions.PriceBelow(
					values.RangeLow(RangeDuration),
				),
			),
		),
		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossFromRange(RangeDuration, 1),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitFromRange(RangeDuration),
			),
		),
		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
