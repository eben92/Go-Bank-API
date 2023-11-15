build:
	go build -o bin/gobank

run: build
	@./bin/gobank

text:
	go test -v ./...