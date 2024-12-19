run:
	@go run cmd/main.go
build:
	@go build -o bin/app cmd/main.go

gen:
	@docker run --rm -v ${shell pwd}:/src -w /src sqlc/sqlc generate
test:
	go test ./...