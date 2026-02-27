.PHONY: build run clean test

build:
	go build -o bin/schedule-bot main.go

run:
	go run main.go

clean:
	rm -rf bin/
	rm -f schedule.db

test:
	go test ./...

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t schedule-bot .

docker-run:
	docker run -d --name schedule-bot schedule-bot