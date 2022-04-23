BINARY=binary
test: 
	go test -v -cover -covermode=atomic ./...

local-build:
	go build -o ${BINARY} cmd/main.go

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

init:
	go mod tidy

run:
	go run cmd/main.go

docker-build:
	docker build -t code-snippets-app .

docker-run:
	docker-compose up --build -d

docker-stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest init build run docker-build docker-run docker-stop vendor lint-prepare lint