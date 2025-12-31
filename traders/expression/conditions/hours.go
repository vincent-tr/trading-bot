package conditions

import (
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func Hours(startHour, endHour int) Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			hour := ctx.Timestamp().Hour()
			return hour >= startHour && hour < endHour
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "Hours",
				formatter.IntValue(startHour),
				formatter.IntValue(endHour),
			)
		},
	)
}
