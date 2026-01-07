package conditions

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

type priceConfig struct {
	offset       int
	candleGetter func(candle *brokers.Candle) float64
	value        values.Value
	comparer     func(price float64, indicatorValue float64) bool
}

func candleOpen(candle *brokers.Candle) float64 {
	return candle.Open
}

func candleHigh(candle *brokers.Candle) float64 {
	return candle.High
}

func candleLow(candle *brokers.Candle) float64 {
	return candle.Low
}

func candleClose(candle *brokers.Candle) float64 {
	return candle.Close
}

func comparerAbove(price float64, indicatorValue float64) bool {
	return price > indicatorValue
}

func comparerBelow(price float64, indicatorValue float64) bool {
	return price < indicatorValue
}

func makePriceConfig(comparer func(price float64, indicatorValue float64) bool, value values.Value, options ...PriceOption) *priceConfig {
	conf := &priceConfig{
		offset:       0,
		comparer:     comparer,
		value:        value,
		candleGetter: candleClose, // Default to close price
	}

	for _, option := range options {
		option.apply(conf)
	}

	return conf
}

func formatPriceOptions(options []PriceOption) []*formatter.FormatterNode {
	nodes := make([]*formatter.FormatterNode, 0, len(options))

	for _, option := range options {
		if option.format != nil {
			nodes = append(nodes, option.format())
		}
	}

	return nodes
}

type PriceOption struct {
	apply  func(*priceConfig)
	format func() *formatter.FormatterNode
}

func Offset(offset int) PriceOption {
	return PriceOption{
		apply: func(ctx *priceConfig) {
			ctx.offset = offset
		},
		format: func() *formatter.FormatterNode {
			return formatter.Function(Package, "Offset", formatter.IntValue(offset))
		},
	}
}

var Open = PriceOption{
	apply: func(ctx *priceConfig) {
		ctx.candleGetter = candleOpen
	},
	format: func() *formatter.FormatterNode {
		return formatter.Value(Package, "Open")
	},
}

var High = PriceOption{
	apply: func(ctx *priceConfig) {
		ctx.candleGetter = candleHigh
	},
	format: func() *formatter.FormatterNode {
		return formatter.Value(Package, "High")
	},
}

var Low = PriceOption{
	apply: func(ctx *priceConfig) {
		ctx.candleGetter = candleLow
	},
	format: func() *formatter.FormatterNode {
		return formatter.Value(Package, "Low")
	},
}

var Close = PriceOption{
	apply: func(ctx *priceConfig) {
		ctx.candleGetter = candleClose
	},
	format: func() *formatter.FormatterNode {
		return formatter.Value(Package, "Close")
	},
}

// PriceAbove returns a condition that checks if the current price is above the given value.
func PriceAbove(value values.Value, options ...PriceOption) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			conf := makePriceConfig(comparerAbove, value, options...)
			return priceCompare(conf, ctx)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceAbove",
				append([]*formatter.FormatterNode{
					value.Format(),
				}, formatPriceOptions(options)...)...,
			)
		},
	)
}

// PriceBelow returns a condition that checks if the current price is below the given value.
func PriceBelow(value values.Value, options ...PriceOption) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			conf := makePriceConfig(comparerBelow, value, options...)
			return priceCompare(conf, ctx)
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"PriceBelow",
				append([]*formatter.FormatterNode{
					value.Format(),
				}, formatPriceOptions(options)...)...,
			)
		},
	)
}

func priceCompare(conf *priceConfig, ctx context.TraderContext) bool {
	candle := ctx.HistoricalData().GetCandle(conf.offset)
	price := conf.candleGetter(&candle)
	indicatorValue := conf.value.Get(ctx)
	return conf.comparer(price, indicatorValue)
}

// ValueAbove returns a condition that checks if valueA is above valueB.
func ValueAbove(valueA, valueB values.Value) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			valueA := valueA.Get(ctx)
			valueB := valueB.Get(ctx)

			return valueA > valueB
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"ValueAbove",
				valueA.Format(),
				valueB.Format(),
			)
		},
	)
}

// ValueBelow returns a condition that checks if valueA is below valueB.
func ValueBelow(valueA, valueB values.Value) Condition {
	return NewCondition(
		func(ctx context.TraderContext) bool {
			valueA := valueA.Get(ctx)
			valueB := valueB.Get(ctx)

			return valueA < valueB
		},
		func() *formatter.FormatterNode {
			return formatter.Function(
				Package,
				"ValueBelow",
				valueA.Format(),
				valueB.Format(),
			)
		},
	)
}
