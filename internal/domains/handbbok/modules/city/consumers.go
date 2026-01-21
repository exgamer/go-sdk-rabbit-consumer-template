package city

import (
	"github.com/exgamer/gosdk-rabbit-core/pkg/config"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city/factories"
)

func GetConsumers(
	consumersFactory *factories.CityConsumersFactory,
) []config.HandlerRegister {
	return []config.HandlerRegister{
		{
			Handler: consumersFactory.CityConsumer.Consume,
			Config: config.NewConsumerTopicDurableConfig(
				"test-consumer",
				"test-rk",
				"test",
				"test",
				100,
			),
		},
	}
}
