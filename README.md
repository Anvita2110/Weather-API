# Weather App Backend (Go + Fiber)

Backend API for fetching current weather data from [WeatherAPI](https://www.weatherapi.com/) with **Redis caching**.

## Features

- **Go Fiber** HTTP server
- **POST** endpoint to fetch current weather by city
- **Redis cache** (10 minute TTL per city)
- **Graceful shutdown** (closes Redis client on SIGINT/SIGTERM)

## Requirements

- Go (see `go.mod` for the version)
- A WeatherAPI key (`WEATHER_API_KEY`)
- Redis (local install or Docker)

## Configuration

The server loads environment variables from a `.env` file if present (optional).

### Environment variables

- **`WEATHER_API_KEY`**: your WeatherAPI key (required)
- **`REDIS_URL`**: Redis connection URL (required)
  - Example (local): `redis://localhost:6379/0`
  - Example (Docker Compose): `redis://redis:6379/0`

## Running locally

1) Start Redis (example with Docker):

```bash
docker run --rm -p 6379:6379 redis:7-alpine
```

2) Create a `.env` file:

```bash
WEATHER_API_KEY=your_key_here
REDIS_URL=redis://localhost:6379/0
```

3) Run the server:

```bash
go run .
```

The server listens on **`http://localhost:3000`**.

## Running with Docker Compose

This repo includes `docker-compose.yml` and a `Dockerfile`.

1) Create a `.env` file (Compose will read it automatically):

```bash
WEATHER_API_KEY=your_key_here
REDIS_PASSWORD=
```

2) Start everything:

```bash
docker compose up --build
```

### Note about Redis env var naming

The app code expects **`REDIS_URL`** (a URL like `redis://redis:6379/0`). If your Compose setup isn’t working, ensure the container gets `REDIS_URL` (the current `docker-compose.yml` uses `REDIS_ADDR`, which the app does not read).

## API

### Health/root

- **GET** `/`
- Response: `Hello, World!`

### Current weather (cached)

- **POST** `/weather/current`
- Body:

```json
{ "city": "London" }
```

- Example:

```bash
curl -sS -X POST "http://localhost:3000/weather/current" \
  -H "Content-Type: application/json" \
  -d '{"city":"London"}'
```

## CORS

`main.go` restricts CORS to `https://vaayu-weather-app.vercel.app/`. For local frontend development, you may need to update `AllowOrigins` (for example, to include `http://localhost:5173`).

