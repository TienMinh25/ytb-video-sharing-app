FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/configs/config.dev.env ./configs/config.dev.env

# rest api port
EXPOSE 3000

# websocket port
EXPOSE 3001

CMD ["./main"]
