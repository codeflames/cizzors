package server

import (
	"fmt"
	"strconv"

	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

type GetCizzorByIDResponse struct {
	Cizzor          models.Cizzor      `json:"cizzor"`
	ClicksAnalytics models.ClickSource `json:"clicks_analytics"`
}

type CreateCizzorRequest struct {
	Url      string `json:"url"`
	Random   bool   `json:"random"`
	ShortUrl string `json:"short_url"`
}

// redirectCizzor redirects the user to the url of the cizzor
// @Summary Redirects the user to the url of the cizzor
// @Description Redirects the user to the url of the cizzor
// @Tags Cizzors
// @Accept json
// @Produce json
// @Param short_url path string true "The short url of the cizzor"
// @Success 307 {string} string	"temporary redirect"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /{short_url} [get]
func redirectCizzor(ctx *fiber.Ctx) error {
	shortUrl := ctx.Params("short_url")

	cizzor, err := models.GetCizzorByShortUrl(shortUrl)

	fmt.Println("initial cizzor: ", cizzor)

	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not get cizzor, " + err.Error(),
			})
	}

	ipAddress := ctx.IP()
	location, err := utils.GetLocation(ipAddress)

	if err != nil {
		fmt.Println("Could not get location: ", err)
	}

	cizzor.UpdateClickSourceCount(ipAddress, location)

	// Increase the count of the cizzor
	cizzor.Count++

	fmt.Println("updated cizzor: ", cizzor)

	cizzor, err = models.UpdateCizzorCount(cizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not update cizzor count",
			})
	}

	return ctx.Redirect(cizzor.Url, fiber.StatusTemporaryRedirect)
}

// gets all the cizzors of the current user
// @Summary Gets all the cizzors of the current user
// @Description Gets all the cizzors of the current user
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} JsonResponse{data=[]models.Cizzor} "All cizzors"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz [get]
func getAllRedirects(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals("user").(models.User)

	cizzors, err := models.GetAllCizzors(currentUser.ID)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(

			ErrorResponse{
				Status:  "error",
				Message: "Could not get cizzors",
			})

	}
	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "All cizzors",
			Data:    cizzors,
		})

}

// generate a qr code for the cizzor
// @Summary Generates a qr code for the cizzor
// @Description Generates a qr code for the cizzor
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Param short_url path string true "The short url of the cizzor"
// @Success 200 {string} string "QR code image"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz/qr/{short_url} [get]
func generateQRCode(ctx *fiber.Ctx) error {
	text := ctx.Params("short_url") // Get the text for the QR code from the query parameter

	//check if it's the owner of the url
	cizzor, err := models.GetCizzorByShortUrl(text)

	currentLocalUser := ctx.Locals("user")

	currentUserId := currentLocalUser.(models.User).ID

	if cizzor.OwnerId != currentUserId {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Unauthorized",
			})
	}

	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not get cizzor, " + err.Error(),
			})
	}

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

// gets a cizzor by id
// @Summary Gets a cizzor by id
// @Description Gets a cizzor by id
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "The id of the cizzor"
// @Success 200 {object} JsonResponse{data=GetCizzorByIDResponse}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz/{id} [get]
func getCizzorById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var err error
	var cizzor models.Cizzor

	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid id",
			})
	}

	currentUser := ctx.Locals("user").(models.User)

	cizzor, err = models.GetCizzorById(currentUser.ID, cizzorId)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not get cizzor " + err.Error(),
			})
	}
	cizzorData, err := models.GetClickSources(currentUser.ID, cizzorId)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not get cizzor data " + err.Error(),
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(

		JsonResponse{
			Status:  "success",
			Message: "Cizzor",
			Data: GetCizzorByIDResponse{
				Cizzor:          cizzor,
				ClicksAnalytics: cizzorData,
			},
		})

}

// creates a cizzor
// @Summary Creates a cizzor
// @Description Creates a cizzor
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Param cizzor body CreateCizzorRequest true "The cizzor to create"
// @Success 200 {object} JsonResponse{data=[]models.Cizzor}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz [post]
func createCizzor(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var request CreateCizzorRequest
	var cizzor models.Cizzor
	var err error

	if err = ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Error parsing json " + err.Error(),
			})
	}
	cizzor.Url = request.Url
	cizzor.Random = request.Random
	cizzor.ShortUrl = request.ShortUrl

	if err = utils.PingURL(cizzor.Url); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Invalid url",
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
			ErrorResponse{
				Status:  "error",
				Message: "Could not create cizzor",
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "Cizzor created",
			Data:    cizzor,
		})
}

type UpdateCizzorRequest struct {
	Url string `json:"url"`
}

// updates a cizzor
// @Summary Updates a cizzor
// @Description Updates a cizzor
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "The id of the cizzor"
// @Param cizzor body UpdateCizzorRequest true "The cizzor to update"
// @Success 200 {object} JsonResponse{data=[]models.Cizzor}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz/{id} [put]
func updateCizzor(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	fmt.Println("cizzor Id " + id)
	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "invalid id",
			})
	}

	request := new(UpdateCizzorRequest)
	// cizzor := new(models.Cizzor)
	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "invalid request body",
			})
	}

	// Create a new user object to store the updated URL field
	// updatedCizzor := models.Cizzor{
	// 	ID:  cizzorId,
	// 	Url: cizzor.Url,
	// }

	currentUser := ctx.Locals("user").(models.User)

	// get the cizzor from the database
	existingCizzor, err := models.GetCizzorById(currentUser.ID, cizzorId)
	if err != nil {
		if err.Error() == "record not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "could not get cizzor: " + err.Error(),
			})
	}

	existingCizzor.Url = request.Url
	// Perform the update
	updatedCizzor, err := models.UpdateCizzor(currentUser.ID, existingCizzor)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "could not update cizzor: " + err.Error(),
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		JsonResponse{
			Status:  "success",
			Message: "cizzor updated successfully",
			Data:    updatedCizzor,
		})
}

// deletes a cizzor
// @Summary Deletes a cizzor
// @Description Deletes a cizzor
// @Tags Cizzors
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "The cizzor id"
// @Success 200 {object} JsonResponse{data=[]models.Cizzor}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cz/{id} [delete]
func deleteCizzor(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var err error

	cizzorId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status:  "error",
				Message: "Could not parse id",
			})
	}
	currentUser := ctx.Locals("user").(models.User)

	err = models.DeleteCizzor(currentUser.ID, cizzorId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			JsonResponse{
				Status:  "error",
				Message: "Could not delete cizzor, " + err.Error(),
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
