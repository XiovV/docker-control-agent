FROM golang:1.17-alpine3.14 AS builder

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go build -o dokkup-agent -ldflags="-s -w"

FROM alpine
WORKDIR /app
COPY --from=builder /app/dokkup-agent /app

CMD ["./dokkup-agent"]
