package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	// "golang.org/x/crypto/bcrypt"
)

// func redirectCizzor(ctx *fiber.Ctx) error {
// 	shortUrl := ctx.Params("short_url")

// 	cizzor, err := models.GetCizzorByShortUrl(shortUrl)
// 	fmt.Println(cizzor)

// 	if err != nil {
// 		if err.Error() == "record not found" {
// 			return ctx.Status(fiber.StatusNotFound).JSON(
// 				JsonResponse{
// 					Status:  "error",
// 					Message: err.Error(),
// 					Data:    []models.Cizzor{},
// 				})
// 		}
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Could not get cizzor",
// 				Data:    []models.Cizzor{},
// 			})
// 	}

// 	// increase the count of the cizzor

// 	fmt.Println(cizzor)

// 	_, err = models.UpdateCizzor(cizzor)

// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Could not update cizzor",
// 				Data:    []models.Cizzor{},
// 			})
// 	}

// 	return ctx.Redirect(cizzor.Url, fiber.StatusTemporaryRedirect)
// }

func redirectCizzor(ctx *fiber.Ctx) error {
	shortUrl := ctx.Params("short_url")

	cizzor, err := models.GetCizzorByShortUrl(shortUrl)

	fmt.Println("initial cizzor: ", cizzor)

	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(JsonResponse{
				Status:  "error",
				Message: err.Error(),
				Data:    []models.Cizzor{},
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
			Status:  "error",
			Message: "Could not get cizzor",
			Data:    []models.Cizzor{},
		})
	}

	// Increase the count of the cizzor
	cizzor.Count++

	fmt.Println("updated cizzor: ", cizzor)

	cizzor, err = models.UpdateCizzor(cizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(JsonResponse{
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

func generateQRCode(ctx *fiber.Ctx) error {
	text := ctx.Params("short_url") // Get the text for the QR code from the query parameter

	fullCizzorUrl := utils.GetFullCizzorUrl(text)

	// Generate the QR code image as a byte slice
	qrCode, err := qrcode.Encode(fullCizzorUrl, qrcode.Medium, 256)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to generate QR code")
	}

	// Set the response content type as image/png
	// ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationPNG)
	ctx.Set(fiber.HeaderContentType, "image/png")

	// Send the QR code image in the response
	return ctx.Send(qrCode)
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

	//get current user and assign it to the cizzor
	user := ctx.Locals("user").(models.User)

	cizzor.OwnerId = user.ID

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

// func updateCizzor(ctx *fiber.Ctx) error {
// 	ctx.Accepts("application/json")

// 	var cizzor models.Cizzor
// 	var err error

// 	if err := ctx.BodyParser(&cizzor); err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Error parsing json " + err.Error(),
// 				Data:    []models.Cizzor{},
// 			})
// 	}

// 	fmt.Println(cizzor)

// 	cizzor, err = models.UpdateCizzor(cizzor)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(
// 			JsonResponse{
// 				Status:  "error",
// 				Message: "Could not update cizzor",
// 				Data:    []models.Cizzor{},
// 			})
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(
// 		JsonResponse{
// 			Status:  "success",
// 			Message: "Cizzor updated",
// 			Data:    cizzor,
// 		})
// }

func updateCizzor(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	fmt.Println("cizzor Id " + id)
	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not parse id",
			})
	}

	cizzor := new(models.Cizzor)
	if err := ctx.BodyParser(cizzor); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
	}

	// Create a new user object to store the updated URL field
	updatedCizzor := models.Cizzor{
		ID:  cizzorId,
		Url: cizzor.Url,
	}

	// Perform the update
	updatedCizzor, err = models.UpdateCizzor(updatedCizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not update cizzor: " + err.Error(),
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "cizzor updated successfully",
			Data: []models.Cizzor{
				{

					ID:       updatedCizzor.ID,
					Url:      updatedCizzor.Url,
					ShortUrl: updatedCizzor.ShortUrl,
					Count:    updatedCizzor.Count,
					Random:   updatedCizzor.Random,
					OwnerId:  updatedCizzor.OwnerId,
				},
			},
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
	fmt.Println(url)
	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to ping URL: %s", err.Error()),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to ping URL: %s", resp.Status),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully pinged URL: %s", url),
	})
}
