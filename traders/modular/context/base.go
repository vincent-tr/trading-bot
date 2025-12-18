package context

import (
	"time"
	"trading-bot/brokers"
	"trading-bot/traders/tools"
)

type TraderContext interface {
	Broker() brokers.Broker
	HistoricalData() *tools.History
	OpenPositions() []brokers.Position
	IndicatorCache() IndicatorCache

	Timestamp() time.Time
	EntryPrice() float64
}

type IndicatorCache interface {
	Tick()
}
