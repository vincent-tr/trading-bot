package indicators

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/values"
)

type cache struct {
	indicators map[string]*Values
}

func NewCache() context.IndicatorCache {
	return &cache{
		indicators: make(map[string]*Values),
	}
}

func (c *cache) Tick() {
	c.indicators = make(map[string]*Values)
}

func (c *cache) access(key string, computer func() []float64) *Values {
	if data, found := c.indicators[key]; found {
		return data
	}
	data := &Values{computer()}
	c.indicators[key] = data
	return data
}

// Indicator is an interface for computing indicator values based on trader context.
type Indicator interface {
	formatter.Formatter
	values.Value
	Values(ctx context.TraderContext) *Values
}

type Values struct {
	data []float64
}

func (values *Values) Current() float64 {
	return values.data[len(values.data)-1]
}

func (values *Values) Previous() float64 {
	return values.data[len(values.data)-2]
}

func (values *Values) All() []float64 {
	return values.data
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

func (i *indicator) Values(ctx context.TraderContext) *Values {
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
