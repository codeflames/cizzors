package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeflames/cizzors/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"

	// "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/codeflames/cizzors/docs"
)

type JsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SetupAndListen() {
	fmt.Println("Setting up server...")

	router := fiber.New()

	cache := cache.New(10*time.Minute, 10*time.Minute)

	// middleware

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Use(logger.New())
	router.Use(compress.New())

	//swagger
	router.Get("/swagger/*", fiberSwagger.WrapHandler) // Path to the generated OpenAPI spec

	// index route
	router.Get("/", func(c *fiber.Ctx) error {
		// Check if the X-Forwarded-For header exists
		if xff := c.Get(fiber.HeaderXForwardedFor); xff != "" {
			// Split the header value to get the client IP address
			ip := strings.Split(xff, ",")[0]
			fmt.Printf("Client IP: %s\n", ip)
		} else {
			// X-Forwarded-For header doesn't exist, use RemoteIP
			ip := c.IP()
			fmt.Printf("Client IP: %s\n", ip)
		}

		return c.Send([]byte("Welcome to Cizzors API"))
	})

	// Create a JWT middleware instance
	jwtMiddleware := middlewares.JwtMiddleware

	// Create a Login middleware instance
	loginMiddleware := middlewares.LoginHashingMiddleware

	// Create a Hashing middleware instance
	hashingMiddleware := middlewares.HashingMiddleware

	// Create a Cache middleware instance
	cacheMiddleware := middlewares.CacheMiddleware(cache)

	invalidateCache := middlewares.InvalidateCacheMiddleware(cache)

	// user routes
	router.Post("/register", hashingMiddleware, createUser)
	router.Post("/login", loginMiddleware, loginUser)
	router.Get("/user/:id", jwtMiddleware, getUserById)
	router.Put("/user/:id", jwtMiddleware, updateUser)
	router.Delete("/user/:id", jwtMiddleware, deleteUser)

	// Create a router group for the routes that require both JWT and Cache middleware
	protectedCizzorRoutes := router.Group("/")
	protectedCizzorRoutes.Use(jwtMiddleware, invalidateCache, cacheMiddleware)

	// Define the routes within the protectedRoutes group
	protectedCizzorRoutes.Get("/cz", getAllRedirects)
	protectedCizzorRoutes.Get("/cz/:id", getCizzorById)
	protectedCizzorRoutes.Post("/cz", createCizzor)
	protectedCizzorRoutes.Put("/cz/:id", updateCizzor)

	// cizzor routes
	router.Get("/:short_url", cacheMiddleware, redirectCizzor)
	router.Get("/cz/qr/:short_url", jwtMiddleware, generateQRCode)
	router.Delete("/cz/:id", jwtMiddleware, deleteCizzor)

	router.Listen(":3001")
}
