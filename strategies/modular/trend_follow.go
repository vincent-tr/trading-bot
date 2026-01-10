package strategies

import (
	"time"
	"trading-bot/common"
	"trading-bot/traders/modular"
	"trading-bot/traders/modular/conditions"
	"trading-bot/traders/modular/indicators"
)

func TrendFollow(strategy modular.StrategyBuilder) {
	strategy.SetFilter(
		conditions.And(
			conditions.HistoryUsable(),
			conditions.NoOpenPositions(),

			conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
			conditions.ExcludeUKHolidays(),
			conditions.ExcludeUSHolidays(),
			conditions.Session(common.LondonSession),
			conditions.Session(common.NYSession),

			conditions.Threshold(indicators.ADX(14), 20.0, conditions.Above),
		),
	)

	strategy.SetLongTrigger(
		conditions.And(
			conditions.IndicatorRange(indicators.RSI(14), 40, 70),
			conditions.PriceThreshold(indicators.EMA(50), conditions.Above),
			conditions.CrossOver(
				indicators.EMA(50),
				indicators.EMA(20),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			conditions.IndicatorRange(indicators.RSI(14), 30, 60),
			conditions.PriceThreshold(indicators.EMA(50), conditions.Below),
			conditions.CrossOver(
				indicators.EMA(50),
				indicators.EMA(20),
				conditions.CrossOverDown,
			),
		),
	)
}
