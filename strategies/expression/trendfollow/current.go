package trendfollow

import (
	"trading-bot/brokers"
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
		expression.HistorySize(600),
		expression.Timeframe(brokers.Timeframe5Minutes),
		expression.Strategy(
			expression.Filter(
				conditions.And(
					conditions.HistoryUsable(),
					conditions.NoOpenPositions(),
				),
			),
			expression.LongTrigger(
				conditions.And(
					// -------------
					// Bullish Trend
					// -------------
					conditions.PriceAbove(
						indicators.EMA(200, indicators.CandleAggregationFactor(3)),
					),
					conditions.ValueAbove(
						indicators.EMA(50, indicators.CandleAggregationFactor(3)),
						indicators.EMA(200, indicators.CandleAggregationFactor(3)),
					),
					// EMA 50 slope is upward
					conditions.ValueAbove(
						// TODO: ATR 5M on 15M candles?
						indicators.NormalizedSlope(
							indicators.EMA(50, indicators.CandleAggregationFactor(3)),
							3,
						),
						values.Constant(0.025),
					),

					// ---------
					// BUY Setup
					// ---------
					conditions.ValueAbove(
						indicators.RSI(14),
						values.Constant(40),
					),
					conditions.PriceBelow(
						indicators.EMA(50),
						conditions.Offset(1),
					),
					conditions.PriceAbove(
						indicators.EMA(50),
					),
				),
			),
			expression.ShortTrigger(
				conditions.And(
					// --------------
					// Bearish Trend
					// --------------
					conditions.PriceBelow(
						indicators.EMA(200, indicators.CandleAggregationFactor(3)),
					),
					conditions.ValueBelow(
						indicators.EMA(50, indicators.CandleAggregationFactor(3)),
						indicators.EMA(200, indicators.CandleAggregationFactor(3)),
					),
					// EMA 50 slope is downward
					conditions.ValueBelow(
						indicators.NormalizedSlope(
							indicators.EMA(50, indicators.CandleAggregationFactor(3)),
							3,
						),
						values.Constant(-0.025),
					),

					// ----------
					// SELL Setup
					// ----------
					conditions.ValueBelow(
						indicators.RSI(14),
						values.Constant(60),
					),
					conditions.PriceAbove(
						indicators.EMA(50),
						conditions.Offset(1),
					),
					conditions.PriceBelow(
						indicators.EMA(50),
					),
				),
			),
		),
		expression.RiskManager(
			expression.StopLoss(
				ordercomputer.StopLossDistance(values.Factor(indicators.ATR(14), 1.2)),
			),
			expression.TakeProfit(
				ordercomputer.TakeProfitRatio(2),
			),
		),
		expression.CapitalAllocator(
			ordercomputer.CapitalFixed(10),
		),
	)
}
