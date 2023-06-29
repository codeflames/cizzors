package server

import (
	"fmt"
	"time"

	"github.com/codeflames/cizzors/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"

	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	cache := cache.New(5*time.Minute, 5*time.Minute)

	limiterConfig := limiter.Config{
		Max:        24,              // Maximum number of requests
		Expiration: 1 * time.Minute, // Time window for rate limiting
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests",
			})
		},
	}

	// middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	router.Use(logger.New())
	router.Use(compress.New())
	router.Use(limiter.New(limiterConfig))

	//swagger
	router.Get("/swagger/*", fiberSwagger.WrapHandler) // Path to the generated OpenAPI spec

	// index route
	// router.Get("/", func(c *fiber.Ctx) error {
	// 	// return c.Send([]byte("Welcome to Cizzors API \n To see the documentation visit: https://cizzors.onrender.com/swagger/index.html"))
	// 	// return html
	// 	return c.SendFile("server/intro.html")

	// })
	router.Static("/static", "./public")

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

	// cizzor routes
	router.Get("/cz/qr/:short_url", jwtMiddleware, generateQRCode)
	router.Delete("/cz/:id", jwtMiddleware, deleteCizzor)
	router.Get("/cz", jwtMiddleware, cacheMiddleware, getAllRedirects)
	router.Get("/:short_url", func(c *fiber.Ctx) error {
		return redirectCizzor(c, cache)
	})
	router.Get("/cz/:id", jwtMiddleware, invalidateCache, getCizzorById)
	router.Post("/cz", jwtMiddleware, invalidateCache, createCizzor)
	router.Put("/cz/:id", updateCizzor)

	// Short URL redirect route

	router.Listen(":3001")
}
