package values

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

const Package string = "values"

// Value is an interface for getting a float64 value.
type Value interface {
	formatter.Formatter
	Get(ctx context.TraderContext) float64
}

type value struct {
	get    func(ctx context.TraderContext) float64
	format func() *formatter.FormatterNode
}

func newValue(
	get func(ctx context.TraderContext) float64,
	format func() *formatter.FormatterNode,
) Value {
	return &value{
		get:    get,
		format: format,
	}
}

func (v *value) Get(ctx context.TraderContext) float64 {
	return v.get(ctx)
}

func (v *value) Format() *formatter.FormatterNode {
	return v.format()
}
