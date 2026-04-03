BINARY=binary

test:
	go test -v -cover -covermode=atomic ./...

test-short:
	go test -short ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

local-build:
	go build -o ${BINARY} ./cmd

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -f coverage.out coverage.html

init:
	go mod tidy

vet:
	go vet ./...

generate:
	npx @redocly/cli bundle openapi/openapi.yaml -o openapi/bundled.yaml
	oapi-codegen --config oapi-codegen-types.yaml openapi/bundled.yaml
	oapi-codegen --config oapi-codegen-server.yaml openapi/bundled.yaml

run:
	go run ./cmd

docker-build:
	docker build -t code-snippets-app .

db-start:
	docker-compose up --build -d

db-stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: test test-short test-coverage local-build clean init vet generate run docker-build db-start db-stop lint-prepare lint
