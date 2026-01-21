package main

import (
	"github.com/exgamer/gosdk-rabbitmq-consumer-template/internal/app"
	"log"
)

func main() {
	appInstance, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = appInstance.RunRabbitKernel()
	if err != nil {
		log.Fatal(err)
	}

	appInstance.WaitForShutdown()
}
