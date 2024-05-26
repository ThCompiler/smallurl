LOG_DIR=./app-log
SWAG_DIRS=./internal/app/delivery/http/v1/,./internal/shortcut/delivery/http/v1/handlers,./internal/shortcut/delivery/http/v1/models/response,./internal/shortcut/delivery/http/v1/models/request,./internal/shortcut/delivery/http/tools

# Генерация

.PHONY: gen-grpc
gen-grpc:
	protoc --go_out=. --go-grpc_out=. ./pkg/grpc/v1/api.proto

# Запуск

.PHONY: run
run:
	sudo docker compose --env-file=./config/env/.env up -d

.PHONY: run-verbose
run-verbose:
	sudo docker compose --env-file=./config/env/.env up

.PHONY: stop
stop:
	sudo docker compose --env-file=./config/env/.env stop

.PHONY: down
down:
	sudo docker compose --env-file=./config/env/.env down

# Сборка

.PHONY: build
build:
	go build -o server -v ./cmd/smallurl

.PHONY: swag-gen
swag-gen:
	swag init --parseDependency --parseInternal --parseDepth 1 -d $(SWAG_DIRS) -g ./swag_info.go -o docs

.PHONY: swag-fmt
swag-fmt:
	swag fmt -d $(SWAG_DIRS) -g ./swag_info.go

.PHONY: build-docker
build-docker:
	sudo docker build --no-cache --network host -f ./docker/smallurl.Dockerfile . --tag smallurl

# Тестирваоние

.PHONY: run-tests
run-tests:
	go test ./...

.PHONY: run-report
run-report:
	allure serve $$(find . -depth -name 'allure-results' -type d | xargs)

.PHONY: clear-reports
clear-reports:
	rm -r $$(find . -depth -name 'allure-results' -type d | xargs)

.PHONY: run-coverage
run-coverage:
	go test -race -covermode=atomic -coverprofile=cover ./...
	cat cover | fgrep -v "mock" | fgrep -v "docs" | fgrep -v "config" | fgrep -v "pb.go"  > cover2
	go tool cover -func=cover2

# Дополнительно

.PHONY: run-hash-test
run-hash-test:
	go run ./hash_test

.PHONY: open-last-log
open-last-log:

.PHONY: open-last-log
open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

.PHONY: clear-logs
clear-logs:
	rm -rf $(LOG_DIR)/*.log

.PHONY: mocks
mocks:
	go generate -n $$(go list ./internal/...)

.PHONY: fmt
fmt:
	gofumpt -e -w -extra .
	goimports -e -w .