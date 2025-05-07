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