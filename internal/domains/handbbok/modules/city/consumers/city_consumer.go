package consumers

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/exgamer/gosdk-rabbit-core/pkg/rabbitmq"
)

func NewCityConsumer(testPublisher *rabbitmq.Publisher) *CityConsumer {
	return &CityConsumer{
		testPublisher: testPublisher,
	}
}

type CityConsumer struct {
	testPublisher *rabbitmq.Publisher
}

func (c *CityConsumer) Consume(ctx context.Context, msg *message.Message) error {
	fmt.Println(string(msg.Payload))

	return nil
}
