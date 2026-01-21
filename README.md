# Шаблон RabbitMQ Listener (gosdk-core + gosdk-rabbit-core)

Этот репозиторий — стартовый шаблон для сервиса‑слушателя RabbitMQ на Go, построенный на:

- `github.com/exgamer/gosdk-core` — приложение, DI, kernel/module менеджеры, конфиги, graceful shutdown
- `github.com/exgamer/gosdk-rabbit-core` — RabbitKernel, publisher/consumer обвязка и конфиги поверх Watermill AMQP

Шаблон показывает:
- как **зарегистрировать и запустить RabbitKernel**
- как **объявлять publishers** декларативно
- как **подключать listeners (consumers/handlers)** через registry
- как организовать код по модулям

---

## Быстрый старт

### 1) Поднять RabbitMQ (локально)

Самый простой способ — Docker:

```bash
docker run -d --name rabbit \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=developer \
  -e RABBITMQ_DEFAULT_PASS=qwerty123 \
  rabbitmq:3-management
```

UI: http://localhost:15672  
Логин/пароль: `developer / qwerty123`

### 2) Настроить окружение

Скопируй пример:

```bash
cp .env.example .env
```

Главные переменные:

- `APP_NAME`, `APP_ENV`, `APP_VERSION`, `DEBUG`
- `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_VHOST`, `RABBITMQ_USER`, `RABBITMQ_PASSWORD`

### 3) Запустить сервис

```bash
go run ./main.go
```

---

## Как устроен проект

Типовая структура (по шаблону):

```
main.go
internal/
  app/                 # bootstrap приложения
  domains/<domain>/    # бизнес-модули
  publishers/          # декларации паблишеров
```

### main.go
Минимальный раннер:
- создаёт app
- запускает rabbit kernel
- ждёт graceful shutdown

---

## RabbitKernel

Kernel подключается и инициализируется в `internal/app` (см. `NewApp()`).

Важно понимать жизненный цикл:

1) `RegisterAndInitKernels(...)` вызывает `RabbitKernel.Init()` — тут регистрируются registry в DI
2) `RegisterAndInitModules(...)` — модули добавляют definitions (listeners/publishers) в registry
3) `RunKernel("rabbit")` вызывает `RabbitKernel.Start()` — тут реально поднимается connection/consumer и т.п.

---

## Listener (consumer): как добавить нового слушателя

В `gosdk-rabbit-core` listener регистрируется как `config.HandlerRegister`:

- `Handler` — функция обработки `(ctx, *message.Message) error`
- `Config` — AMQP конфиг (exchange/queue/bind/qos/consumer tag)

### 1) Создай consumer (обработчик)

Пример (как в шаблоне `CityConsumer`):

```go
type CityConsumer struct {}

func (c *CityConsumer) Consume(ctx context.Context, msg *message.Message) error {
    fmt.Println(string(msg.Payload))
    return nil
}
```

### 2) Собери список handlers модуля

Пример (как `internal/domains/.../consumers.go`):

```go
return []config.HandlerRegister{
  {
    Handler: consumersFactory.CityConsumer.Consume,
    Config: config.NewConsumerTopicDurableConfig(
      "test-consumer", // consumer tag (идентификатор консьюмера)
      "test-rk",       // routing key
      "test",          // exchange
      "test",          // queue
      100,             // prefetch
    ),
  },
}
```

> **Consumer tag** в вашей версии задаётся через поле `ConsumeConfig.Consumer` внутри `NewConsumer*Config`.
> Это имя будет видно в RabbitMQ UI в списке consumers у очереди.

### 3) Зарегистрируй handlers в registry в Init() модуля

Пример логики модуля (псевдо‑код, адаптируй под свой модуль):

```go
func (m *Module) Init(a *app.App) error {
    reg, err := di.GetRabbitConsumersRegistry(a.Container)
    if err != nil { return err }

    consumersFactory := factories.NewCityConsumersFactory(testPublisher)
    handlers := GetConsumers(consumersFactory)

    reg.RegisterMultipleHandler(handlers)
    return nil
}
```

Kernel сам подхватит все зарегистрированные handlers при `RunKernel("rabbit")`.

---

## Publisher: как объявить и использовать

### 1) Объяви publisher definition

Шаблон уже содержит `internal/publishers/publishers.go`:

```go
func GetPublishers() []config.PublisherDefinition {
  return []config.PublisherDefinition{
    {
      Name:   "TestPublisher",
      Config: config.NewPublisherTopicDurableConfig("test"),
    },
  }
}
```

Где:
- `Name` — имя паблишера в registry
- `NewPublisherTopicDurableConfig("test")` — exchange = `test`

### 2) Зарегистрируй publishers в registry (в Init() модуля)

```go
pubReg, err := di.GetRabbitPublishersRegistry(a.Container)
if err != nil { return err }

_ = pubReg.RegisterMultiple(publishers.GetPublishers())
```

### 3) Получи publisher по имени и публикуй

```go
testPublisher, err := pubReg.Get("TestPublisher")
if err != nil { return err }

payload := map[string]any{"id": 1, "name": "hello"}

err = testPublisher.Publish("test-rk", payload) // topic = routing key
if err != nil { return err }
```

С метаданными:

```go
err = testPublisher.PublishWithMetaData(
  "test-rk",
  map[string]string{
    "trace_id": "abc-123",
  },
  payload,
)
```

---

## Рекомендации по именованию

Routing keys:
- `domain.event`, пример: `orders.paid`, `products.updated`

Consumer tag:
- `<service>-<env>-<instance>`, пример: `catalog-local-macbook`

Queue:
- лучше уникальная на сервис/модуль (чтобы не мешать разные сервисы)

---

## Troubleshooting

### “Kernel не запускается”
Убедись, что вызывается именно:
- `RunKernel("rabbit")`
  а не только `RegisterAndInitKernels(...)`.

В шаблоне это сделано в `RunRabbitKernel()`.

### “Не вижу consumer в RabbitMQ UI”
Проверь:
- сервис реально запущен (`go run main.go`)
- правильные `RABBITMQ_*` переменные
- queue/exchange/routingKey совпадают с тем, куда ты публикуешь

---

## Лицензия / заметки
Шаблон предназначен для внутреннего использования в проектах на базе `gosdk-core` и `gosdk-rabbit-core`.
