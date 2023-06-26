package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codeflames/cizzors/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func createUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Invalid request body",
				Data:    []models.User{},
			})
	} else {

		// Retrieve the hashed password from the request context using ctx.Locals()
		hashedPassword, ok := ctx.Locals("hashedPassword").(string)
		if !ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
				Status:  "error",
				Message: "An error occurred while creating user",
				Data:    []models.User{},
			})
		}

		// Set the hashed password in the user object
		user.Password = hashedPassword

		user, err := models.CreateUser(*user)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) || strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_users_username"`) {
				return ctx.Status(fiber.StatusConflict).JSON(
					JsonResponse{
						Status:  "error",
						Message: "User already exists",
						Data:    []models.User{},
					})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				JsonResponse{
					Status:  "error",
					Message: "Could not create user",
					Data:    []models.User{},
				})
		}
		return ctx.Status(fiber.StatusCreated).JSON(
			JsonResponse{
				Status:  "success",
				Message: "User created successfully, please login to continue",
				// strip password from response
				Data: []models.User{
					{
						ID:       user.ID,
						Username: user.Username,
						Email:    user.Email,
					},
				},
			})
	}
}

func getUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	//parse the id to a uint64 and return an error if it fails
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
			})
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not get user",
				Data:    []models.User{},
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "User retrieved successfully",
			Data: []models.User{
				{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
				},
			},
		})
}

func loginUser(ctx *fiber.Ctx) error {
	// user := new(models.LoginUserModel)
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(JsonResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	responseUser, err := models.LoginUser(*user)
	if err != nil {
		// no user found
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(JsonResponse{
				Status:  "error",
				Message: "User not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
			Status:  "error",
			Message: "Error: " + err.Error(),
		})
	}

	userId := strconv.FormatUint(responseUser.ID, 10)

	// Generate a JWT token
	token, err := generateToken(
		userId,
		responseUser.Username,
		responseUser.Email,
		os.Getenv("JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
			Status:  "error",
			Message: "Error generating JWT token",
		})
	}

	// Return the response with the token
	return ctx.Status(fiber.StatusOK).JSON(JsonResponse{
		Status:  "success",
		Message: "User logged in successfully",
		Data: map[string]interface{}{
			"user": []models.User{
				{
					ID:       responseUser.ID,
					Email:    responseUser.Email,
					Username: responseUser.Username,
				},
			},
			"token": token,
		},
	})
}

// func updateUser(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	fmt.Println("cizzor Id "id)
// 	cizzorId, err := strconv.ParseUint(id, 10, 64)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Could not parse id",
// 			})
// 	}

// 	cizzor := new(models.Cizzor)
// 	if err := ctx.BodyParser(cizzor); err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Invalid request body",
// 			})
// 	}

// 	// Create a new user object to store the updated URL field
// 	updatedCizzor := models.Cizzor{
// 		ID:  cizzorId,
// 		Url: cizzor.Url,
// 	}

// 	// Perform the update
// 	updatedCizzor, err = models.UpdateCizzor(updatedCizzor)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Could not update cizzor: " + err.Error(),
// 			})
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(
// 		JsonResponse{
// 			Status:  "success",
// 			Message: "cizzor updated successfully",
// 			Data: []models.Cizzor{
// 				{

// 					ID:       updatedCizzor.ID,
// 					Url:      updatedCizzor.Url,
// 					ShortUrl: updatedCizzor.ShortUrl,
// 					Count:    updatedCizzor.Count,
// 					Random:   updatedCizzor.Random,
// 					OwnerId:  updatedCizzor.OwnerId,
// 				},
// 			},
// 		})
// }

func updateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
				// Data:    []models.User{},
			})
	}

	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Invalid request body",
				// Data:    []models.User{},
			})
	} else {
		user.ID = userId
		user, err := models.UpdateUser(*user)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				JsonResponse{
					Status:  "error",
					Message: "Could not update user: " + err.Error(),
					// Data:    []models.User{},
				})
		}
		return ctx.Status(fiber.StatusOK).JSON(
			JsonResponse{
				Status:  "success",
				Message: "User updated successfully",
				Data: []models.User{
					{
						ID:       user.ID,
						Username: user.Username,
						Email:    user.Email,
					},
				},
			})
	}
}

func deleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	//parse the id to a uint64 and return an error if it fails
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
			})
	}

	err = models.DeleteUser(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not delete user" + err.Error(),
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "User deleted successfully",
		})
}

func generateToken(userID, username, email string, jwtSecret string) (string, error) {
	// Create the claims for the token
	claims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	}

	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the JWT secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
