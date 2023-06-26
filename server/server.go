package server

import (
	"fmt"

	"github.com/codeflames/cizzors/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "golang.org/x/crypto/bcrypt"
)

type JsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SetupAndListen() {
	fmt.Println("Setting up server...")

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Use(logger.New())
	router.Use(compress.New())
	router.Use(cache.New())

	// user routes
	router.Post("/register", middlewares.HashingMiddleware, createUser)
	router.Post("/login", middlewares.LoginHashingMiddleware, loginUser)
	router.Get("/user/:id", middlewares.JwtMiddleware, getUserById)
	router.Put("/user/:id", middlewares.JwtMiddleware, updateUser)
	router.Delete("/user/:id", middlewares.JwtMiddleware, deleteUser)

	// cizzor routes
	router.Get("/cz", middlewares.JwtMiddleware, getAllRedirects)
	router.Get("/:short_url", redirectCizzor)
	router.Get("/qrcode/generate/:short_url", generateQRCode)
	router.Get("/cz/:id", middlewares.JwtMiddleware, getCizzorById)
	router.Post("/cz", middlewares.JwtMiddleware, createCizzor)
	router.Put("/cz/:id", middlewares.JwtMiddleware, updateCizzor)

	router.Delete("/cz/:id", middlewares.JwtMiddleware, deleteCizzor)

	router.Listen(":3001")
}
