package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
	"weather-app/cache"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
)

type WeatherRequest struct {
	City string `json:"city"`
}

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
		AirQuality struct {
			Pm2_5 float64 `json:"pm2_5"`
			Pm10  float64 `json:"pm10"`
		} `json:"air_quality"`
	} `json:"current"`
}

var redisClient *cache.RedisClient

func CloseRedisClient() {
	if redisClient != nil {
		redisClient.Close()
	}
}

func SetupRoutes(app *fiber.App) {
	API_KEY := os.Getenv("WEATHER_API_KEY")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_DB := 0

	if redisClient == nil {
		redisClient = cache.NewRedisClient(REDIS_ADDR, REDIS_PASSWORD, REDIS_DB)
	}

	weatherRoute := app.Group("/weather")
	weatherRoute.Post("/current", func(c fiber.Ctx) error {
		body := c.Body()

		var weatherRequest WeatherRequest

		if err := json.Unmarshal(body, &weatherRequest); err != nil {
			return c.Status(400).SendString("Invalid JSON format")
		}

		city := weatherRequest.City
		if city == "" {
			return c.Status(400).SendString("City field is required")
		}

		cacheKey := fmt.Sprintf("weather:%s", city)
		var cachedResponse WeatherResponse

		if err := redisClient.Get(c.Context(), cacheKey, &cachedResponse); err == nil {
			log.Printf("Cache hit for city: %s", city)
			return c.JSON(cachedResponse)
		}

		baseURL := "http://api.weatherapi.com/v1/current.json"
		params := url.Values{}
		params.Add("key", API_KEY)
		params.Add("q", city)
		params.Add("aqi", "yes")

		fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		cli := client.New()
		resp, err := cli.Get(fullURL)
		if err != nil {
			return c.Status(500).SendString("Failed to fetch weather data")
		}
		defer resp.Close()

		var weatherResponse WeatherResponse

		err = json.Unmarshal(resp.Body(), &weatherResponse)
		if err != nil {
			return c.Status(500).SendString("Failed to parse weather data")
		}

		if err := redisClient.Set(c.Context(), cacheKey, weatherResponse, 10*time.Minute); err != nil {
			log.Printf("Failed to cache response: %v", err)
		} else {
			log.Printf("Cached weather data for city: %s", city)
		}

		return c.JSON(weatherResponse)
	})
}
