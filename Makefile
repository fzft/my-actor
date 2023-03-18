build:
	go build -o bin/actor


run: build
	./bin/actor


test:
	go test ./...