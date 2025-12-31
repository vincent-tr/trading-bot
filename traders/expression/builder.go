package expression

import (
	"fmt"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/ordercomputer"
)

type Configuration struct {
	historySizeConfiguration
	strategyConfiguration
	riskManagerConfiguration
	capitalAllocatorConfiguration
}

func Builder(
	historySize *historySizeConfiguration,
	strategy *strategyConfiguration,
	riskManager *riskManagerConfiguration,
	capitalAllocator *capitalAllocatorConfiguration,
) *Configuration {
	return &Configuration{
		historySizeConfiguration:      *historySize,
		strategyConfiguration:         *strategy,
		riskManagerConfiguration:      *riskManager,
		capitalAllocatorConfiguration: *capitalAllocator,
	}
}

type historySizeConfiguration struct {
	historySize int
}

func HistorySize(size int) *historySizeConfiguration {
	return &historySizeConfiguration{historySize: size}
}

type strategyConfiguration struct {
	filter       conditions.Condition
	longTrigger  conditions.Condition
	shortTrigger conditions.Condition
}

func Strategy(filter *strategyFilterConfiguration, longTrigger *strategyLongTriggerConfiguration, shortTrigger *strategyShortTriggerConfiguration) *strategyConfiguration {
	return &strategyConfiguration{
		filter:       filter.value,
		longTrigger:  longTrigger.value,
		shortTrigger: shortTrigger.value,
	}
}

type strategyFilterConfiguration struct {
	value conditions.Condition
}

func Filter(value conditions.Condition) *strategyFilterConfiguration {
	return &strategyFilterConfiguration{value}
}

type strategyLongTriggerConfiguration struct {
	value conditions.Condition
}

func LongTrigger(value conditions.Condition) *strategyLongTriggerConfiguration {
	return &strategyLongTriggerConfiguration{value}
}

type strategyShortTriggerConfiguration struct {
	value conditions.Condition
}

func ShortTrigger(value conditions.Condition) *strategyShortTriggerConfiguration {
	return &strategyShortTriggerConfiguration{value}
}

type riskManagerConfiguration struct {
	stopLoss   ordercomputer.OrderComputer
	takeProfit ordercomputer.OrderComputer
}

func RiskManager(stopLoss *riskManagerStopLossConfiguration, takeProfit *riskManagerTakeProfitConfiguration) *riskManagerConfiguration {
	return &riskManagerConfiguration{
		stopLoss:   stopLoss.value,
		takeProfit: takeProfit.value,
	}
}

type riskManagerStopLossConfiguration struct {
	value ordercomputer.OrderComputer
}

func StopLoss(value ordercomputer.OrderComputer) *riskManagerStopLossConfiguration {
	return &riskManagerStopLossConfiguration{value}
}

type riskManagerTakeProfitConfiguration struct {
	value ordercomputer.OrderComputer
}

func TakeProfit(value ordercomputer.OrderComputer) *riskManagerTakeProfitConfiguration {
	return &riskManagerTakeProfitConfiguration{value}
}

type capitalAllocatorConfiguration struct {
	allocator ordercomputer.OrderComputer
}

func CapitalAllocator(value ordercomputer.OrderComputer) *capitalAllocatorConfiguration {
	return &capitalAllocatorConfiguration{value}
}

func (config *Configuration) Format() *formatter.FormatterNode {
	return formatter.Format("ModularTrader",
		formatter.Format(fmt.Sprintf("HistorySize: %d", config.historySize)),
		formatter.FormatWithChildren("Filter", config.filter),
		formatter.FormatWithChildren("LongTrigger", config.longTrigger),
		formatter.FormatWithChildren("ShortTrigger", config.shortTrigger),
		formatter.FormatWithChildren("StopLoss", b.stopLoss),
		formatter.FormatWithChildren("TakeProfit", b.takeProfit),
		formatter.FormatWithChildren("CapitalAllocator", b.capitalAllocator),
	)
}

func Format(b *Configuration) string {
	bu, ok := b.(*builder)
	if !ok {
		panic(fmt.Sprintf("invalid builder type: %T", b))
	}

	return bu.Format().Detailed()
}
