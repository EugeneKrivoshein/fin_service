FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go mod tidy
RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]