package middlewares

import (
	"fmt"

	"github.com/codeflames/cizzors/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Define hashing middleware
func HashingMiddleware(c *fiber.Ctx) error {
	// Example: Hashing password from the request body

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	fmt.Println(user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to hash password")
	}

	fmt.Println("hashed password: ", string(hashedPassword))

	// Replace the original password with the hashed password in the request context
	// c.Context().SetUserValue("hashedPassword", string(hashedPassword))
	c.Locals("hashedPassword", string(hashedPassword))

	return c.Next()
}

func LoginHashingMiddleware(c *fiber.Ctx) error {
	// Example: Hashing password from the request body

	// user := models.LoginUserModel{}
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// Replace the original password with the hashed password in the request context

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to hash password")
	}
	c.Context().SetUserValue("hashedPassword", string(hashedPassword))

	return c.Next()
}
