package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

// / OrderComputer is an interface for computing orders properties based on trader context.
type OrderComputer interface {
	formatter.Formatter
	Compute(ctx context.TraderContext, order *brokers.Order) error
}

func newOrderComputer(
	compute func(ctx context.TraderContext, order *brokers.Order) error,
	format func() *formatter.FormatterNode,
) OrderComputer {
	return &orderComputer{
		compute: compute,
		format:  format,
	}
}

type orderComputer struct {
	compute func(ctx context.TraderContext, order *brokers.Order) error
	format  func() *formatter.FormatterNode
}

func (oc *orderComputer) Compute(ctx context.TraderContext, order *brokers.Order) error {
	return oc.compute(ctx, order)
}

func (oc *orderComputer) Format() *formatter.FormatterNode {
	return oc.format()
}
