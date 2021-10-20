FROM golang:1.17-alpine3.14 AS builder

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go build -o docker_control_agent -ldflags="-s -w"

FROM alpine
WORKDIR /app
COPY --from=builder /app/docker_control_agent /app

CMD ["./docker_control_agent"]
