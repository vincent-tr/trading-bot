package conditions

import (
	"time"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Weekdays(weekdays ...time.Weekday) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			for _, day := range weekdays {
				if ctx.Timestamp().Weekday() == day {
					return true
				}
			}
			return false
		},
		func() *formatter.FormatterNode {
			values := make([]*formatter.FormatterNode, len(weekdays))
			for i, day := range weekdays {
				values[i] = formatter.Value("time", day.String())
			}
			return formatter.Function(Package, "Weekdays", values...)
		},
	)
}
