package expression

import (
	"fmt"
	"maps"
	"slices"
	"time"
	"trading-bot/brokers"
	"trading-bot/common"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/context"
	"trading-bot/traders/expression/indicators"
	"trading-bot/traders/expression/ordercomputer"
	"trading-bot/traders/tools"
)

var log = common.NewLogger("traders/expression")

func Setup(broker brokers.Broker, config *Configuration) error {

	trader, err := newTrader(broker, config)
	if err != nil {
		return err
	}

	log.Debug("%s", config.Format().Detailed())

	broker.RegisterMarketDataCallback(brokers.Timeframe1Minute, func(candle brokers.Candle) {
		trader.tick(candle)
	})

	return nil
}

type trader struct {
	broker           brokers.Broker
	history          *tools.History
	openPositions    map[brokers.Position]struct{}
	indicatorCache   context.IndicatorCache
	filter           conditions.Condition
	longTrigger      conditions.Condition
	shortTrigger     conditions.Condition
	stopLoss         ordercomputer.OrderComputer
	takeProfit       ordercomputer.OrderComputer
	capitalAllocator ordercomputer.OrderComputer
}

func newTrader(broker brokers.Broker, config *Configuration) (*trader, error) {

	if config.historySize <= 0 {
		return nil, fmt.Errorf("history size must be greater than 0")
	}
	if config.filter == nil {
		return nil, fmt.Errorf("filter must be set")
	}
	if config.longTrigger == nil && config.shortTrigger == nil {
		return nil, fmt.Errorf("either long or short trigger must be set")
	}
	if config.stopLoss == nil {
		return nil, fmt.Errorf("stop loss computer must be set")
	}
	if config.takeProfit == nil {
		return nil, fmt.Errorf("take profit computer must be set")
	}
	if config.capitalAllocator == nil {
		return nil, fmt.Errorf("capital allocator must be set")
	}

	return &trader{
		broker:           broker,
		history:          tools.NewHistory(config.historySize),
		openPositions:    make(map[brokers.Position]struct{}),
		indicatorCache:   indicators.NewCache(),
		filter:           config.filter,
		longTrigger:      config.longTrigger,
		shortTrigger:     config.shortTrigger,
		stopLoss:         config.stopLoss,
		takeProfit:       config.takeProfit,
		capitalAllocator: config.capitalAllocator,
	}, nil
}

func (t *trader) tick(candle brokers.Candle) {
	t.history.AddCandle(candle)

	for pos := range t.openPositions {
		if pos.Closed() || pos.Canceled() {
			delete(t.openPositions, pos)
		}
	}

	t.indicatorCache.Tick()

	if !t.filter.Execute(t) {
		return
	}

	shouldTakeLong := t.longTrigger.Execute(t)
	shouldTakeShort := t.shortTrigger.Execute(t)

	if shouldTakeLong && shouldTakeShort {
		log.Warning("Both long and short triggers are true, ignoring")
		return
	}

	if shouldTakeLong {
		t.takePosition(brokers.PositionDirectionLong)
	}

	if shouldTakeShort {
		t.takePosition(brokers.PositionDirectionShort)
	}
}

func (t *trader) takePosition(direction brokers.PositionDirection) {
	order := &brokers.Order{
		Direction: direction,
	}

	err := t.stopLoss.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute stop loss: %v", err)
		return
	}

	err = t.takeProfit.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute take profit: %v", err)
		return
	}

	err = t.capitalAllocator.Compute(t, order)
	if err != nil {
		log.Error("Failed to compute capital allocation: %v", err)
		return
	}

	// TODO: reason

	pos, err := t.broker.PlaceOrder(order)
	if err != nil {
		log.Error("Failed to place order: %v", err)
		return
	}

	t.openPositions[pos] = struct{}{}
}

var _ context.TraderContext = (*trader)(nil)

func (t *trader) Broker() brokers.Broker {
	return t.broker
}
func (t *trader) HistoricalData() *tools.History {
	return t.history
}

func (t *trader) OpenPositions() []brokers.Position {
	return slices.Collect(maps.Keys(t.openPositions))
}

func (t *trader) IndicatorCache() context.IndicatorCache {
	return t.indicatorCache
}

func (t *trader) Timestamp() time.Time {
	return t.broker.GetCurrentTime()
}

func (t *trader) EntryPrice() float64 {
	return t.history.GetPrice()
}
