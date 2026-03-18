FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o weather-app

FROM alpine:latest

# Install redis-cli for debugging (optional)
RUN apk add --no-cache redis

WORKDIR /app

COPY --from=builder /app/weather-app .

EXPOSE 3000

CMD ["./weather-app"]