# wb.tech internship l0 task
Простой микросервис на стеке технологий golang, postgres и nats-streaming

Реализованы:
1. In-Memory кэширование, кэш потокобезопасный!
2. Graceful shutdown
3. Migrations

## Как собрать и поднять
1. Клонировать репо
2. Внутри директории проекта запустить
   ```
   ~$ make build
   ~$ make run
   ``` 
Сервис слушает по адресу localhost:3000
Мониторинг nats-streaming по адресу - localhost:6222
 
Методы API:

GET - /order={id} - получить заказ по ID
GET - /time - получить время на сервере (тестовый ендпоинт)
