package city

import (
	"github.com/exgamer/gosdk-core/pkg/app"
	"github.com/exgamer/gosdk-rabbit-core/pkg/di"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city/dto"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city/factories"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/publishers"
)

// Module модуль бизнес логики
type Module struct {
}

func (m *Module) Name() string {
	return "city"
}

func (m *Module) Init(a *app.App) error {
	//регистрируем паблишеры
	pubRegistry, err := di.GetRabbitPublishersRegistry(a.Container)
	if err != nil {
		return err
	}

	pubs := publishers.GetPublishers()

	err = pubRegistry.RegisterMultiple(pubs)
	if err != nil {
		return err
	}

	// Пример получения паблишера и отправка сообщения в тестовую очередь
	testPublisher, err := pubRegistry.Get("TestPublisher")
	if err != nil {
		return err
	}

	payload := dto.TestPayload{
		ID:    1,
		Name:  "Test user",
		Event: "created",
	}
	err = testPublisher.Publish("test-rk", payload)
	if err != nil {
		return err
	}

	//регистрируем консьюмеры
	consumersFactory := factories.NewCityConsumersFactory(testPublisher)

	consumers := GetConsumers(consumersFactory)

	reg, err := di.GetRabbitConsumersRegistry(a.Container) // твой helper для DI
	if err != nil {
		return err
	}

	reg.RegisterMultipleHandler(consumers)

	return nil
}
