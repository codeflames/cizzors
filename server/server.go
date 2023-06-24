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

func redirectCizzor(ctx *fiber.Ctx) error {
	shortUrl := ctx.Params("short_url")

	cizzor, err := models.GetCizzorByShortUrl(shortUrl)
	fmt.Println(cizzor)

	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				JsonResponse{
					Status:  "error",
					Message: err.Error(),
					Data:    []models.Cizzor{},
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not get cizzor",
				Data:    []models.Cizzor{},
			})
	}

	// increase the count of the cizzor

	cizzor.Count = cizzor.Count + 1

	fmt.Println(cizzor)

	_, err = models.UpdateCizzor(cizzor)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not update cizzor",
				Data:    []models.Cizzor{},
			})
	}

	return ctx.Redirect(cizzor.Url, fiber.StatusTemporaryRedirect)
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
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
				Data:    []models.Cizzor{},
			})
	}

	_, err = models.UpdateCizzor(cizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not update cizzor",
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

	if err = ctx.BodyParser(&cizzor); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Error parsing json " + err.Error(),
				Data:    []models.Cizzor{},
			})
	}
	if err = pingURL(ctx, cizzor.Url); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Invalid url",
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

func updateCizzor(ctx *fiber.Ctx) error {
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

	cizzor, err = models.UpdateCizzor(cizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not update cizzor",
				Data:    []models.Cizzor{},
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "Cizzor updated",
			Data:    cizzor,
		})
}

func deleteCizzor(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var err error

	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
				Data:    []models.Cizzor{},
			})
	}

	err = models.DeleteCizzor(cizzorId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not delete cizzor",
				Data:    []models.Cizzor{},
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "Cizzor deleted",
			Data:    []models.Cizzor{},
		})
}

func pingURL(ctx *fiber.Ctx, url string) error {
	resp := ctx.Get(url)

	fmt.Println(resp)

	return nil
}

// if err != nil {
// 	return err
// }
// defer resp.Body()

// if resp.StatusCode() != fiber.StatusOK {
// 	return fmt.Errorf("received non-200 status code: %d", resp.StatusCode())
// }

// return nil
// }

func SetupAndListen() {
	fmt.Println("Setting up server...")

	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Get("/cz", getAllRedirects)
	router.Get("/cz/:id", getCizzorById)
	router.Post("/cz", createCizzor)
	router.Put("/cz/:id", updateCizzor)
	router.Get("/:shorturl", redirectCizzor)
	router.Delete("/cz/:id", deleteCizzor)

	router.Listen(":3001")
}
