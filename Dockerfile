FROM golang:1.17-alpine3.14 AS builder

RUN apk add --no-cache bash
RUN apk add --no-cache upx
WORKDIR /app

COPY . .

EXPOSE 8080

RUN go build -o docker_control_agent -ldflags="-s -w"
RUN upx -9 -k docker_control_agent

FROM alpine
WORKDIR /app
COPY --from=builder /app/docker_control_agent /app

CMD ["./docker_control_agent"]
