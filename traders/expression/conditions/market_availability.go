package conditions

import (
	"time"
	"trading-bot/common"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

func HistoryUsable() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return ctx.HistoricalData().IsUsable()
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "HistoryUsable")
		},
	)
}

func NoOpenPositions() Condition {
	return newCondition(
		func(ctx context.TraderContext) bool {
			return len(ctx.OpenPositions()) == 0
		},
		func() *formatter.FormatterNode {
			return formatter.Function(Package, "NoOpenPositions")
		},
	)
}

type SessionType int

const (
	SessionLondon SessionType = iota
	SessionNewYork
)

func Session(sessionType SessionType) Condition {
	var session *common.Session

	switch sessionType {
	case SessionLondon:
		session = common.LondonSession
	case SessionNewYork:
		session = common.NYSession
	default:
		panic("unknown session type")
	}

	return newCondition(
		func(ctx context.TraderContext) bool {
			return session.IsOpen(ctx.Timestamp())
		},
		func() *formatter.FormatterNode {
			var value *formatter.FormatterNode

			switch sessionType {
			case SessionLondon:
				value = formatter.Value(Package, "SessionLondon")
			case SessionNewYork:
				value = formatter.Value(Package, "SessionNewYork")
			default:
				panic("unknown session type")
			}

			return formatter.Function(Package, "Session", value)
		},
	)
}

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
