package middlewares

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func CacheMiddleware(cache *cache.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			// Only cache GET requests
			return c.Next()
		}

		cacheKey := GenerateCacheKey(c)

		// Check if the response is already in the cache
		if cached, found := cache.Get(cacheKey); found {
			fmt.Println("cache hit")
			c.Response().Header.Set("X-Cache-Status", "hit")
			// Set custom header for cache hit
			return c.JSON(cached)
		}

		err := c.Next()
		if err != nil {
			return err
		}

		// Check if the response is an error
		if c.Response().StatusCode() >= 400 {
			fmt.Println("cache skip due to error response")
			return nil
		}

		var data interface{}
		body := c.Response().Body()

		if len(body) > 0 {
			// Unmarshal the response body
			err = json.Unmarshal(body, &data)
			if err != nil {
				return c.JSON(fiber.Map{"error": err.Error()})
			}
		}

		// Cache the response for 10 minutes
		cache.Set(cacheKey, data, 10*time.Minute)

		fmt.Println("cache miss")
		c.Response().Header.Set("X-Cache-Status", "miss") // Set custom header for cache miss

		return nil
	}
}

func GenerateCacheKey(c *fiber.Ctx) string {
	// Generate a cache key from the request path and query parameters
	cacheKey := c.Path()

	// check to see if query params has :short_url
	if c.Params("short_url") != "" {

		// Exclude query parameters from cache key for PUT and POST requests
		if c.Method() != fiber.MethodGet {
			cacheKey = strings.Split(cacheKey, "?")[0]
		}

		return cacheKey
	} else {
		return cacheKey + "?" + c.Params("short_url")
	}
}

// Create a custom middleware to invalidate cache for PUT and POST requests
func InvalidateCacheMiddleware(cache *cache.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Proceed with the request
		err := c.Next()
		if err != nil {
			return err
		}

		// Invalidate the cache for PUT and POST requests
		if c.Method() == fiber.MethodPut || c.Method() == fiber.MethodPost {
			cache.Flush()
		}

		return nil
	}
}
