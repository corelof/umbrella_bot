include .env
export

build:
	go build
run: build
	./umbrella_bot
test: 
	go test ./...