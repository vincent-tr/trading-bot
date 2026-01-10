package strategies

import (
	"time"
	"trading-bot/common"
	"trading-bot/traders/modular"
	"trading-bot/traders/modular/conditions"
	"trading-bot/traders/modular/indicators"
)

func Reversal(strategy modular.StrategyBuilder) {
	strategy.SetFilter(
		conditions.And(
			conditions.HistoryUsable(),
			conditions.NoOpenPositions(),
			conditions.Weekday(time.Tuesday, time.Wednesday, time.Thursday),
			conditions.Session(common.LondonSession),
			conditions.Session(common.NYSession),
			conditions.ExcludeUKHolidays(),
			conditions.ExcludeUSHolidays(),
		),
	)

	strategy.SetLongTrigger(
		conditions.And(
			conditions.PriceThreshold(indicators.EMA(50), conditions.Above),
			conditions.CrossOver(
				indicators.Const(14, 30.0),
				indicators.RSI(14),
				conditions.CrossOverUp,
			),
		),
	)

	strategy.SetShortTrigger(
		conditions.And(
			conditions.PriceThreshold(indicators.EMA(50), conditions.Below),
			conditions.CrossOver(
				indicators.Const(14, 70.0),
				indicators.RSI(14),
				conditions.CrossOverDown,
			),
		),
	)
}
