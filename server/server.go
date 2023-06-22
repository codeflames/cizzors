package server

import (
	"fmt"
	"strconv"

	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type JsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func getAllRedirects(ctx *fiber.Ctx) error {
	cizzors, err := models.GetAllCizzors()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not get cizzors",
				Data:    []models.Cizzor{},
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "All cizzors",
			Data:    cizzors,
		})

}

func getCizzorById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var err error
	var cizzor models.Cizzor

	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
				Data:    []models.Cizzor{},
			})
	}

	cizzor, err = models.GetCizzorById(cizzorId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not get cizzor " + err.Error(),
				Data:    []models.Cizzor{},
			})
	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "Success",
			Data:    cizzor,
		})

}

func createCizzor(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var cizzor models.Cizzor
	var err error

	if err := ctx.BodyParser(&cizzor); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Error parsing json " + err.Error(),
				Data:    []models.Cizzor{},
			})
	}
	if cizzor.Random {
		cizzor.ShortUrl = utils.RandomUrl()
	}

	cizzor, err = models.CreateCizzor(cizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not create cizzor",
				Data:    []models.Cizzor{},
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "Cizzor created",
			Data:    cizzor,
		})
}

func SetupAndListen() {
	fmt.Println("Setting up server...")

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Get("/", getAllRedirects)
	router.Get("/:id", getCizzorById)
	router.Post("/", createCizzor)

	router.Listen(":3001")
}
