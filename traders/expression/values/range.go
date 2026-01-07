package values

import (
	"math"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

type rangeConfig struct {
	period int
	offset int
}

func makeRangeConfig(period int, options []RangeOption) *rangeConfig {
	conf := &rangeConfig{
		period: period,
		offset: 1, // Default offset is 1 to exclude the current candle
	}

	for _, option := range options {
		option.apply(conf)
	}

	return conf
}

func formatRangeOptions(options []RangeOption) []*formatter.FormatterNode {
	nodes := make([]*formatter.FormatterNode, 0, len(options))

	for _, option := range options {
		if option.format != nil {
			nodes = append(nodes, option.format())
		}
	}

	return nodes
}

type RangeOption struct {
	apply  func(*rangeConfig)
	format func() *formatter.FormatterNode
}

func Offset(offset int) RangeOption {
	return RangeOption{
		apply: func(ctx *rangeConfig) {
			ctx.offset = offset
		},
		format: func() *formatter.FormatterNode {
			return formatter.Function(Package, "Offset", formatter.IntValue(offset))
		},
	}
}

// RangeSize computes the range (high - low) over a specified number of historical candles.
// It returns the difference between the highest high and lowest low within the period.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeSize(20, 1) calculates the range from candles 1-20 bars ago (skipping the most recent completed candle).
func RangeSize(period int, options ...RangeOption) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			conf := makeRangeConfig(period, options)

			high := getRangeHigh(ctx, conf)
			low := getRangeLow(ctx, conf)
			return high - low
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeSize",
				append([]*formatter.FormatterNode{
					formatter.IntValue(period),
				}, formatRangeOptions(options)...)...,
			)
		},
	)
}

// RangeHigh returns the highest high price over a specified number of historical candles.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeHigh(20, 1) finds the highest price from candles 1-20 bars ago.
func RangeHigh(period int, options ...RangeOption) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			conf := makeRangeConfig(period, options)
			return getRangeHigh(ctx, conf)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeHigh",
				append([]*formatter.FormatterNode{
					formatter.IntValue(period),
				}, formatRangeOptions(options)...)...,
			)
		},
	)
}

// RangeLow returns the lowest low price over a specified number of historical candles.
// The offset parameter shifts the range window backward in time (offset=1 excludes only the last candle).
// For example, RangeLow(20, 1) finds the lowest price from candles 1-20 bars ago.
func RangeLow(period int, options ...RangeOption) Value {
	return NewValue(
		func(ctx context.TraderContext) float64 {
			conf := makeRangeConfig(period, options)
			return getRangeLow(ctx, conf)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"RangeLow",
				append([]*formatter.FormatterNode{
					formatter.IntValue(period),
				}, formatRangeOptions(options)...)...,
			)
		},
	)
}

func getRangeHigh(ctx context.TraderContext, conf *rangeConfig) float64 {
	highs := ctx.HistoricalData().GetHighPrices()

	hightest := math.Inf(-1)

	for i := conf.offset; i < conf.period+conf.offset; i++ {
		price := highs.At(i)
		if price > hightest {
			hightest = price
		}
	}

	return hightest
}

func getRangeLow(ctx context.TraderContext, conf *rangeConfig) float64 {
	lows := ctx.HistoricalData().GetLowPrices()

	lowest := math.Inf(1)

	for i := conf.offset; i < conf.period+conf.offset; i++ {
		price := lows.At(i)
		if price < lowest {
			lowest = price
		}
	}

	return lowest
}
