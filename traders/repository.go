package traders

import (
	"trading-bot/brokers"
	"trading-bot/traders/basic"
	"trading-bot/traders/expression"
	"trading-bot/traders/gpt"
	"trading-bot/traders/modular"
)

type GptConfig = gpt.Config

func SetupBasicTrader(broker brokers.Broker) {
	basic.Setup(broker)
}

func SetupGptTrader(broker brokers.Broker, config *GptConfig) {
	gpt.Setup(broker, config)
}

func SetupModularTrader(broker brokers.Broker, builder modular.Builder) error {
	return modular.Setup(broker, builder)
}

func SetupExpessionTrader(broker brokers.Broker, config *expression.Configuration) error {
	return expression.Setup(broker, config)
}
