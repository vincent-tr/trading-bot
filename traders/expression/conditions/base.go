package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

const Package string = "conditions"

// / Condition represents a trading condition that can be evaluated within a trading context.
type Condition interface {
	formatter.Formatter
	Execute(ctx context.TraderContext) bool
}

func newCondition(
	execute func(ctx context.TraderContext) bool,
	format func() *formatter.FormatterNode,
) Condition {
	return &condition{
		execute: execute,
		format:  format,
	}
}

type condition struct {
	execute func(ctx context.TraderContext) bool
	format  func() *formatter.FormatterNode
}

func (c *condition) Execute(ctx context.TraderContext) bool {
	return c.execute(ctx)
}

func (c *condition) Format() *formatter.FormatterNode {
	return c.format()
}
