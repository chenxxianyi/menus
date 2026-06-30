.PHONY: build run test clean docker-up docker-down

build:
	go build -o bin/server.exe ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development with hot reload (requires air)
dev:
	air -c .air.toml
