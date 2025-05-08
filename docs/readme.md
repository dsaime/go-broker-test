Сервис для расчета профита по данным трейдов.
Состоит из двух исполняемых файлов, работающих с одной базой данных.

Перед запуском сервиса, необходимо создать базу данных и воссоздать схему из `./migrations/`.
Для миграции можно применять утилиту [migrate](https://github.com/golang-migrate/migrate)

Пример:
```sh
go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ./migrations/sqlite -database sqlite3://data1.db up # Восстановить схему из скриптов миграции
```

Запуск: 

Запустить api, позволяющее работать с данными. 
```sh
go run ./cmd/server.go --db data.db --listen 8080
```

Запустить worker, который будет работать в фоне и обновлять статистику по аккаунту. 
```sh
go run ./cmd/worker.go --db data.db --poll 100ms
``` 

Вспомогательные команды:

Стандартный линтер:
```sh
go vet ./...
```
Стандартный + golangci-lint:
```sh
go vet ./... && golangci-lint run -v -j $(( $(nproc) - 1))
```
Запуск тестов: 
```sh
go test -race ./...
```

Погенерить трейды:
```sh
for i in {1..10000}; do curl -X POST --location "http://localhost:8080/trades" -H "Content-Type: application/json" -d '{
  "account": "qweuietwk",
  "symbol": "EURUSD",
  "volume": 9.0,
  "open": 3339.0,
  "close": 3399.0,
  "side": "buy"
}' & sleep 0.01; done; wait; echo AllOk!
```

Посмотреть распределение по воркерам:
```sql
select worker_id, job_status, count(1)
from trades_q
group by worker_id, job_status
```