package factories

import (
	"github.com/exgamer/gosdk-rabbit-core/pkg/rabbitmq"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city/consumers"
)

func NewCityConsumersFactory(testPublisher *rabbitmq.Publisher) *CityConsumersFactory {
	return &CityConsumersFactory{
		CityConsumer: consumers.NewCityConsumer(testPublisher),
	}
}

type CityConsumersFactory struct {
	CityConsumer *consumers.CityConsumer
}
