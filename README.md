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

## Для более детализированной информации ознакомьтесь в документацей к следующим пакетам

- [Документация GOSDK-CORE](https://github.com/exgamer/gosdk-core)

- [Документация GOSDK-RABBIT-CORE](https://github.com/exgamer/gosdk-rabbit-core?tab=readme-ov-file)

---

## 🚀 Что даёт шаблон

- 🧠 Единый `App` lifecycle
- 🌐 Готовый HTTP kernel (Gin)
- 🧩 Dependency Injection из коробки
- 🧱 Модульная архитектура (business modules)
- ⚙️ Конфигурация через env
- ♻️ Graceful shutdown
- ❗ Стандартизированная обработка ошибок
- 🧪 Удобная база для тестирования

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
