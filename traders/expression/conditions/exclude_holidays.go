package conditions

import (
	"trading-bot/common"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

type Holidays int

const (
	UKHolidays = iota
	USHolidays
)

func ExcludeHolidays(typ Holidays) Condition {
	switch typ {
	case UKHolidays:
		return newCondition(
			func(ctx context.TraderContext) bool {
				return !common.IsUKHoliday(ctx.Timestamp())
			},
			func() *formatter.FormatterNode {
				value := formatter.Value(Package, "UKHolidays")
				return formatter.Function(Package, "ExcludeHolidays", value)
			},
		)

	case USHolidays:
		return newCondition(
			func(ctx context.TraderContext) bool {
				return !common.IsUSHoliday(ctx.Timestamp())
			},
			func() *formatter.FormatterNode {
				value := formatter.Value(Package, "USHolidays")
				return formatter.Function(Package, "ExcludeHolidays", value)
			},
		)

	default:
		panic("unknown holiday type")
	}
}
