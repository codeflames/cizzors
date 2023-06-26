package middlewares

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codeflames/cizzors/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(c *fiber.Ctx) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing authorization token")
	}

	// Check if the authorization header starts with "Bearer "
	if len(authHeader) > 7 && strings.ToLower(authHeader[0:7]) == "bearer " {
		tokenString := authHeader[7:] // Remove "Bearer " prefix from the token string

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token signing method")
			}

			// Provide the secret key used for signing the token
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid authorization token")
		}

		// Verify token claims and extract user information
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract the user information from the claims

			// check the run time type of the claims["userID"] interface{} value
			// if it is not a string, return an error

			userID := claims["userID"].(string)

			// convert userID to uint64
			userIDUint, err := strconv.ParseUint(userID, 10, 64)
			if err != nil {
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).SendString("Invalid authorization token")
			}

			user := models.User{
				ID:       userIDUint,
				Username: claims["username"].(string),
				Email:    claims["email"].(string),
				// Add other user fields as needed
			}

			// Set the user information in the request context for future handlers
			c.Locals("user", user)

			return c.Next()
		}
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid authorization token")
}
