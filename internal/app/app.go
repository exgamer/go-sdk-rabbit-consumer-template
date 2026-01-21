package app

import (
	"github.com/exgamer/gosdk-core/pkg/app"
	app2 "github.com/exgamer/gosdk-rabbit-core/pkg/app"
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/domains/handbbok/modules/city"
)

type App struct {
	*app.App
}

func NewApp() (*App, error) {
	appInstance := &App{
		App: app.NewApp(),
	}

	err := appInstance.RegisterAndInitKernels(
		app2.NewRabbitKernel().EnableConsumer().EnablePublisher(),
	)
	if err != nil {
		return nil, err
	}

	err = appInstance.RegisterAndInitModules(
		&city.Module{},
	)
	if err != nil {
		return nil, err
	}

	return appInstance, nil
}

func (a *App) RunRabbitKernel() error {
	if err := a.RunKernel(app2.RabbitKernelName); err != nil {
		return err
	}

	return nil
}
