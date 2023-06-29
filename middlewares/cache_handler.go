// package middlewares

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/patrickmn/go-cache"
// )

// type JsonResponse struct {
// 	Status  string      `json:"status"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data"`
// }

// func CacheMiddleware(cache *cache.Cache) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		if c.Method() != "GET" {
// 			// Only cache GET requests
// 			return c.Next()
// 		}

// 		cacheKey := c.Path() + "?" + c.Params("id") // Generate a cache key from the request path and query parameters

// 		// Check if the response is already in the cache
// 		if cached, found := cache.Get(cacheKey); found {
// 			fmt.Println("cache hit")
// 			return c.JSON(cached)
// 		}
// 		err := c.Next()
// 		if err != nil {
// 			return err
// 		}

// 		var data JsonResponse
// 		cacheKey = c.Path() + "?" + c.Params("id")

// 		body := c.Response().Body()
// 		err = json.Unmarshal(body, &data)
// 		if err != nil {
// 			return c.JSON(fiber.Map{"error": err.Error()})
// 		}

// 		// Cache the response for 10 minutes
// 		cache.Set(cacheKey, data, 10*time.Minute)

// 		fmt.Println("cache miss")

// 		return nil
// 	}
// }

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

		cacheKey := generateCacheKey(c)

		// Check if the response is already in the cache
		if cached, found := cache.Get(cacheKey); found {
			fmt.Println("cache hit")
			c.Response().Header.Set("X-Cache-Status", "hit") // Set custom header for cache hit
			return c.JSON(cached)
		}

		err := c.Next()
		if err != nil {
			return err
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

func generateCacheKey(c *fiber.Ctx) string {
	// Generate a cache key from the request path and query parameters
	cacheKey := c.Path()

	// Exclude query parameters from cache key for PUT and POST requests
	if c.Method() != fiber.MethodGet {
		cacheKey = strings.Split(cacheKey, "?")[0]
	}

	return cacheKey
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
