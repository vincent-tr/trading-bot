package conditions

import (
	"trading-bot/common"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/formatter"
)

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
