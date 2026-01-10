package strategies

import (
	"trading-bot/common"
	"trading-bot/traders/modular"
	"trading-bot/traders/modular/conditions"
	"trading-bot/traders/modular/indicators"
)

func Simple(strategy modular.StrategyBuilder) {

	strategy.SetFilter(conditions.And(
		conditions.HistoryUsable(),
		conditions.NoOpenPositions(),

		// conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
		// conditions.ExcludeUKHolidays(),
		// conditions.ExcludeUSHolidays(),
		conditions.Session(common.LondonSession),
		conditions.Session(common.NYSession),

		conditions.Compare(indicators.ATR(14), indicators.Mean(indicators.ATR(14), 10), conditions.Above),
	))

	strategy.SetLongTrigger(
		conditions.And(
			conditions.Slope(indicators.EMA(200), 10, conditions.SlopeRising),
			conditions.PriceThreshold(indicators.EMA(200), conditions.Above),
			conditions.PriceThreshold(indicators.EMA(20), conditions.Below),
			conditions.DistancePriceToPrevious(indicators.EMA(20), 1, conditions.Above),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			conditions.Slope(indicators.EMA(200), 10, conditions.SlopeFalling),
			conditions.PriceThreshold(indicators.EMA(200), conditions.Below),
			conditions.PriceThreshold(indicators.EMA(20), conditions.Above),
			conditions.DistancePriceToPrevious(indicators.EMA(20), 1, conditions.Below),
		),
	)
}
