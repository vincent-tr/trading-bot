package ordercomputer

import (
	"trading-bot/brokers"
	"trading-bot/traders/modular/context"
	"trading-bot/traders/modular/formatter"
	"trading-bot/traders/modular/marshal"
)

type OrderComputer interface {
	formatter.Formatter
	Compute(ctx context.TraderContext, order *brokers.Order) error
	ToJsonSpec() (string, any)
}

func newOrderComputer(
	compute func(ctx context.TraderContext, order *brokers.Order) error,
	format func() *formatter.FormatterNode,
	toJsonSpec func() (string, any),
) OrderComputer {
	return &orderComputer{
		compute:    compute,
		format:     format,
		toJsonSpec: toJsonSpec,
	}
}

type orderComputer struct {
	compute    func(ctx context.TraderContext, order *brokers.Order) error
	format     func() *formatter.FormatterNode
	toJsonSpec func() (string, any)
}

func (oc *orderComputer) Compute(ctx context.TraderContext, order *brokers.Order) error {
	return oc.compute(ctx, order)
}

func (oc *orderComputer) Format() *formatter.FormatterNode {
	return oc.format()
}

func (oc *orderComputer) ToJsonSpec() (string, any) {
	return oc.toJsonSpec()
}

var jsonParsers = marshal.NewRegistry[OrderComputer]()

func FromJSON(jsonData []byte) (OrderComputer, error) {
	return jsonParsers.FromJSON(jsonData)
}

func ToJSON(oc OrderComputer) ([]byte, error) {
	panic("ToJSON not implemented for order computers")
}
