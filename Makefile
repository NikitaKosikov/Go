.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

debug: build
	docker-compose up --remove-orphans debug

test:
	go test --short  ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out

export DB_URI=mongodb://localhost:27019
export DB_NAME=testDb
export CONTAINER_NAME=test_db

test.integration:
	docker run --rm -d  -p 27019:27017 --name test_db -e MONGODB_DATABASE=testDb mongo:4.2.23-bionic
	go test -v ./tests/
	docker stop test_db