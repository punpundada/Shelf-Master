FROM golang:1.23.1-alpine3.20 AS base
# ARG CGO_ENABLED=0
WORKDIR /usr/src/app



FROM base AS development
RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN source .env && go build -o bin/app cmd/main.go
CMD ["air"]


FROM scratch
COPY --from=development /usr/src/app/cmd/app /usr/src/app/cmd/app
ENTRYPOINT ["./app"]