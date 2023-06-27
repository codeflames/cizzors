package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/utils"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
}

type CreateAccountRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
}

// createUser creates a new user.
// @Summary Create user
// @Description Create a new user.
// @Tags User
// @Accept json
// @Produce json
// @Param user body CreateAccountRequest{} true "CreateAccountRequest"
// @Success 201 {object} JsonResponse{data=UserResponse} "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Could not create user"
// @Router /register [post]
func createUser(ctx *fiber.Ctx) error {
	userRequest := new(CreateAccountRequest)
	user := new(models.User)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
	} else {

		err = utils.EmailVerifier(userRequest.Email)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Invalid email address",
				})
		}

		if !utils.PasswordLengthChecker(userRequest.Password) {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Password must be at least 6 characters long",
				})
		}

		if !utils.UsernameLengthChecker(userRequest.Username) {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Username must be at least 3 characters long",
				})
		}

		user.Password = userRequest.Password
		user.Username = userRequest.Username
		user.Email = userRequest.Email

		// Retrieve the hashed password from the request context using ctx.Locals()
		hashedPassword, ok := ctx.Locals("hashedPassword").(string)
		if !ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "An error occurred while creating the user"})
		}

		// Set the hashed password in the user object
		user.Password = hashedPassword

		user, err := models.CreateUser(*user)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) || strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_users_username"`) {
				return ctx.Status(fiber.StatusConflict).JSON(
					ErrorResponse{
						Status:  "error",
						Message: "User already exists",
					})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Could not create user",
				})
		}

		return ctx.Status(fiber.StatusCreated).JSON(
			JsonResponse{
				Status:  "success",
				Message: "User created successfully, please login to continue",
				// strip password from response
				Data: UserResponse{
					ID:        user.ID,
					Username:  user.Username,
					Email:     user.Email,
					CreatedAt: user.CreatedAt.Format(time.RFC3339),
				},
			})
	}
}

// getUserById retrieves a user by id.
// @Summary Get user by id
// @Description Get a user by id.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security Bearer
// @Success 200 {object} JsonResponse{data=UserResponse} "User retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Could not get user"
// @Router /user/{id} [get]
func getUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	//parse the id to a uint64 and return an error if it fails
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid user id",
			})
	}

	// check if the current user is the same as the user being requested
	// if not, return an error
	currentLocalUser := ctx.Locals("user")

	currentUserId := currentLocalUser.(models.User).ID

	if currentUserId != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "You are not authorized to view this user",
			})
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "User not found",
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not get user",
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "User retrieved successfully",
			Data: UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.Format(time.RFC3339),
			},
		})
}

// loginUser logs in a user.
// @Summary Login user
// @Description Login a user.
// @Tags User
// @Accept json
// @Produce json
// @Param user body LoginRequest{} true "LoginRequest"
// @Success 200 {object} JsonResponse{data=LoginResponse} "User logged in successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Error logging in user"
// @Router /login [post]
func loginUser(ctx *fiber.Ctx) error {
	// user := new(models.LoginUserModel)
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
	}

	// Validate the user input
	err := utils.EmailVerifier(user.Email)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid email address",
			})
	}

	responseUser, err := models.LoginUser(*user)
	if err != nil {
		// no user found
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "User not found",
				})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Error, " + err.Error(),
			})
	}

	userId := strconv.FormatUint(responseUser.ID, 10)

	// Generate a JWT token
	token, err := utils.GenerateToken(
		userId,
		responseUser.Username,
		responseUser.Email,
		os.Getenv("JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Error, " + err.Error(),
			})
	}

	// Return the response with the token
	return ctx.Status(fiber.StatusOK).JSON(JsonResponse{
		Status:  "success",
		Message: "User logged in successfully",
		Data: LoginResponse{
			Token:     token,
			Username:  responseUser.Username,
			Email:     responseUser.Email,
			ID:        responseUser.ID,
			CreatedAt: responseUser.CreatedAt.Format(time.RFC3339),
		},
	})
}

// updateUser updates a user.
// @Summary Update user
// @Description Update a user.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UserUpdateRequest{} true "UserUpdateRequest"
// @Security Bearer
// @Success 200 {object} JsonResponse{data=UserResponse} "User updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Could not update user"
// @Router /user/{id} [put]
func updateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	// user := new(models.User)
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid user id",
			})
	}

	userRequest := new(UserUpdateRequest)
	if err := ctx.BodyParser(userRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
	} else {
		// check if the current user is the same as the user being requested
		// if not, return an error
		currentLocalUser := ctx.Locals("user")

		currentUserId := currentLocalUser.(models.User).ID

		if currentUserId != userId {
			return ctx.Status(fiber.StatusForbidden).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "You are not authorized to update this user",
				})
		}
		// get the user from the database
		existingUser, err := models.GetUserByID(userId)
		if err != nil {
			if err.Error() == "record not found" {
				return ctx.Status(fiber.StatusNotFound).JSON(
					ErrorResponse{
						Status:  "error",
						Message: "User not found",
					})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Could not get user, " + err.Error(),
				})
		}

		existingUser.Username = userRequest.Username

		user, err := models.UpdateUser(existingUser)
		if err != nil {
			if err.Error() == "record not found" {
				return ctx.Status(fiber.StatusNotFound).JSON(
					ErrorResponse{
						Status:  "error",
						Message: "User not found",
					})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "Could not update user, " + err.Error(),
				})
		}
		return ctx.Status(fiber.StatusOK).JSON(
			JsonResponse{
				Status:  "success",
				Message: "User updated successfully",
				Data: UserResponse{
					ID:        user.ID,
					Username:  user.Username,
					Email:     user.Email,
					CreatedAt: user.CreatedAt.Format(time.RFC3339),
				},
			})
	}
}

// deleteUser deletes a user.
// @Summary Delete user
// @Description Delete a user.
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security Bearer
// @Success 200 {object} JsonResponse{data=UserResponse} "User deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid user id"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Could not delete user"
// @Router /user/{id} [delete]
func deleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	//parse the id to a uint64 and return an error if it fails
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid user id",
			})
	}
	// check if the current user is the same as the user being requested
	// if not, return an error
	currentLocalUser := ctx.Locals("user")

	currentUserId := currentLocalUser.(models.User).ID

	if currentUserId != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "You are not authorized to delete this user",
			})
	}

	err = models.DeleteUser(userId)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: "User not found",
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not delete user, " + err.Error(),
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "User deleted successfully",
		})
}
