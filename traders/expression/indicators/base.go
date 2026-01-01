package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
	"trading-bot/traders/tools"
)

const Package string = "indicators"

type cache struct {
	indicators map[string]*tools.Values
}

func NewCache() context.IndicatorCache {
	return &cache{
		indicators: make(map[string]*tools.Values),
	}
}

func (c *cache) Tick() {
	c.indicators = make(map[string]*tools.Values)
}

func (c *cache) access(key string, computer func() []float64) *tools.Values {
	if data, found := c.indicators[key]; found {
		return data
	}
	data := tools.NewValues(computer())
	c.indicators[key] = data
	return data
}

// Indicator is an interface for computing indicator values based on trader context.
type Indicator interface {
	formatter.Formatter
	values.Value
	Values(ctx context.TraderContext) *tools.Values
}

type indicator struct {
	compute func(ctx context.TraderContext) []float64
	format  func() *formatter.FormatterNode
}

func newIndicator(
	compute func(ctx context.TraderContext) []float64,
	format func() *formatter.FormatterNode,
) Indicator {
	return &indicator{
		compute: compute,
		format:  format,
	}
}

func (i *indicator) Values(ctx context.TraderContext) *tools.Values {
	c := ctx.IndicatorCache().(*cache)
	key := i.format().Compact()

	return c.access(key, func() []float64 {
		return i.compute(ctx)
	})
}

func (i *indicator) Get(ctx context.TraderContext) float64 {
	return i.Values(ctx).Current()
}

func (i *indicator) Format() *formatter.FormatterNode {
	return i.format()
}
