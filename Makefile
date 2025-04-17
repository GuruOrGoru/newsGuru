run: build
	@./bin/server-api

build:
	@go build -o bin/server-api cmd/server-api/main.go

clean:
	@rm -f bin/server-api
test:
	@go test -v ./pkg/tests
