package expression

import (
	"trading-bot/brokers"
	"trading-bot/traders/expression/conditions"
	"trading-bot/traders/expression/formatter"
	"trading-bot/traders/expression/ordercomputer"
)

const Package string = "expression"

type Configuration struct {
	historySizeConfiguration
	timeframeConfiguration
	strategyConfiguration
	riskManagerConfiguration
	capitalAllocatorConfiguration
}

func (config *Configuration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"Builder",
		config.historySizeConfiguration.Format(),
		config.timeframeConfiguration.Format(),
		config.strategyConfiguration.Format(),
		config.riskManagerConfiguration.Format(),
		config.capitalAllocatorConfiguration.Format(),
	)
}

func Builder(
	historySize *historySizeConfiguration,
	timeframe *timeframeConfiguration,
	strategy *strategyConfiguration,
	riskManager *riskManagerConfiguration,
	capitalAllocator *capitalAllocatorConfiguration,
) *Configuration {
	return &Configuration{
		historySizeConfiguration:      *historySize,
		timeframeConfiguration:        *timeframe,
		strategyConfiguration:         *strategy,
		riskManagerConfiguration:      *riskManager,
		capitalAllocatorConfiguration: *capitalAllocator,
	}
}

type historySizeConfiguration struct {
	historySize int
}

func (config *historySizeConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"HistorySize",
		formatter.IntValue(config.historySize),
	)
}

func HistorySize(size int) *historySizeConfiguration {
	return &historySizeConfiguration{historySize: size}
}

type timeframeConfiguration struct {
	timeframe brokers.Timeframe
}

func (config *timeframeConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"Timeframe",
		formatter.Value("brokers", config.timeframe.Format()),
	)
}

func Timeframe(timeframe brokers.Timeframe) *timeframeConfiguration {
	return &timeframeConfiguration{timeframe}
}

type strategyConfiguration struct {
	filter       *strategyFilterConfiguration
	longTrigger  *strategyLongTriggerConfiguration
	shortTrigger *strategyShortTriggerConfiguration
}

func (config *strategyConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"Strategy",
		config.filter.Format(),
		config.longTrigger.Format(),
		config.shortTrigger.Format(),
	)
}

func Strategy(filter *strategyFilterConfiguration, longTrigger *strategyLongTriggerConfiguration, shortTrigger *strategyShortTriggerConfiguration) *strategyConfiguration {
	return &strategyConfiguration{
		filter,
		longTrigger,
		shortTrigger,
	}
}

type strategyFilterConfiguration struct {
	value conditions.Condition
}

func (config *strategyFilterConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"Filter",
		config.value.Format(),
	)
}

func Filter(value conditions.Condition) *strategyFilterConfiguration {
	return &strategyFilterConfiguration{value}
}

type strategyLongTriggerConfiguration struct {
	value conditions.Condition
}

func (config *strategyLongTriggerConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"LongTrigger",
		config.value.Format(),
	)
}

func LongTrigger(value conditions.Condition) *strategyLongTriggerConfiguration {
	return &strategyLongTriggerConfiguration{value}
}

type strategyShortTriggerConfiguration struct {
	value conditions.Condition
}

func (config *strategyShortTriggerConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"ShortTrigger",
		config.value.Format(),
	)
}

func ShortTrigger(value conditions.Condition) *strategyShortTriggerConfiguration {
	return &strategyShortTriggerConfiguration{value}
}

type riskManagerConfiguration struct {
	stopLoss   *riskManagerStopLossConfiguration
	takeProfit *riskManagerTakeProfitConfiguration
}

func (config *riskManagerConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"RiskManager",
		config.stopLoss.Format(),
		config.takeProfit.Format(),
	)
}

func RiskManager(stopLoss *riskManagerStopLossConfiguration, takeProfit *riskManagerTakeProfitConfiguration) *riskManagerConfiguration {
	return &riskManagerConfiguration{
		stopLoss,
		takeProfit,
	}
}

type riskManagerStopLossConfiguration struct {
	value ordercomputer.OrderComputer
}

func (config *riskManagerStopLossConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"StopLoss",
		config.value.Format(),
	)
}

func StopLoss(value ordercomputer.OrderComputer) *riskManagerStopLossConfiguration {
	return &riskManagerStopLossConfiguration{value}
}

type riskManagerTakeProfitConfiguration struct {
	value ordercomputer.OrderComputer
}

func (config *riskManagerTakeProfitConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"TakeProfit",
		config.value.Format(),
	)
}

func TakeProfit(value ordercomputer.OrderComputer) *riskManagerTakeProfitConfiguration {
	return &riskManagerTakeProfitConfiguration{value}
}

type capitalAllocatorConfiguration struct {
	capitalAllocator ordercomputer.OrderComputer
}

func (config *capitalAllocatorConfiguration) Format() *formatter.FormatterNode {
	return formatter.Function(
		Package,
		"CapitalAllocator",
		config.capitalAllocator.Format(),
	)
}

func CapitalAllocator(value ordercomputer.OrderComputer) *capitalAllocatorConfiguration {
	return &capitalAllocatorConfiguration{value}
}
