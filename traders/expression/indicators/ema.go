package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"

	"github.com/markcheno/go-talib"
)

type emaConfig struct {
	period                 int
	candleAggregtionFactor int
}

func makeEmaConfig(period int, options []EmaOption) *emaConfig {
	conf := &emaConfig{
		period:                 period,
		candleAggregtionFactor: 1, // Default is no aggregation
	}

	for _, option := range options {
		option.apply(conf)
	}

	return conf
}

type EmaOption struct {
	apply  func(*emaConfig)
	format func() *formatter.FormatterNode
}

func CandleAggregationFactor(factor int) EmaOption {
	return EmaOption{
		apply: func(ctx *emaConfig) {
			ctx.candleAggregtionFactor = factor
		},
		format: func() *formatter.FormatterNode {
			return formatter.Function(Package, "CandleAggregationFactor", formatter.IntValue(factor))
		},
	}
}

// EMA computes the Exponential Moving Average over the specified period.
// It uses the closing prices of the candles for the calculation.
//
// When CandleAggregationFactor option is provided, it aggregates the candles by the specified factor before calculating the EMA.
// For example, if CandleAggregationFactor(3) is used, every 3 consecutive candles are averaged into one before computing the EMA.
// So this allows using a higher timeframe EMA on lower timeframe data. (15min EMA on 5min data, etc.)
//
// Note: when using aggregation
// - ensure that the historical data has enough candles to accommodate the aggregation factor.
// - do not use previous candle offsets that are not multiples of the aggregation factor, as this may lead to unexpected results.
func EMA(period int, options ...EmaOption) Indicator {
	return NewIndicator(
		func(ctx context.TraderContext) []float64 {
			config := makeEmaConfig(period, options)

			closePrices := ctx.HistoricalData().GetClosePrices().All()
			closePrices = aggregateClosePrices(closePrices, config.candleAggregtionFactor)
			return talib.Ema(closePrices, period)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "EMA", formatter.IntValue(period))
		},
	)
}

func aggregateClosePrices(prices []float64, factor int) []float64 {
	if factor <= 1 {
		return prices
	}

	aggregated := make([]float64, 0, len(prices)/factor)
	for i := 0; i < len(prices); i += factor {
		aggregated = append(aggregated, prices[i])
	}
	return aggregated
}
