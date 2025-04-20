package rest

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *bannerHandler) createBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	description := c.FormValue("desc")
	if description == "" {
		log.Error().Msg("handler::createBanner - description is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Description is required"))
	}

	if len(description) > constants.MaxLengthDescriptionInput {
		log.Error().Msg("handler::createBanner - description exceeds max length")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Description exceeds max length"))
	}

	isXssSafe := utils.IsXSSSafe(description)
	if !isXssSafe {
		log.Error().Msg("handler::createBanner - description contains XSS")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Description is invalid"))
	}

	mf, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("handler::createBanner - Failed to parse multipart form")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid multipart form data"))
	}

	fileHeaders := mf.File[constants.MultipartFormFile]
	switch len(fileHeaders) {
	case 0:
		log.Error().Msg("handler::createBanner - No image file uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file is required"))
	case 1:
		log.Info().Msgf("handler::createBanner - %s file is valid", fileHeaders[0].Filename)
	default:
		log.Error().Msgf("handler::createBanner - too many files uploaded: %d", len(fileHeaders))
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Only one image file is allowed"))
	}

	fh := fileHeaders[0]
	if fh.Size > constants.MaxFileSize {
		log.Error().Msg("handler::createBanner - File size exceeds limit")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !constants.AllowedImageExtensions[ext] {
		log.Error().Msg("handler::createBanner - Invalid file extension")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
	}

	file, err := fh.Open()
	if err != nil {
		log.Error().Err(err).Msg("handler::createBanner - Failed to open file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("handler::createBanner - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
	}

	mimeType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		log.Error().Msg("handler::createBanner - Uploaded file is not a valid image")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
	}

	objectName := "banner_images/" + utils.GenerateBucketFileUUID() + ext
	uploadFile := dto.UploadFileRequest{
		ObjectName:     objectName,
		File:           fileBytes,
		FileHeaderSize: fh.Size,
		ContentType:    mimeType,
		Filename:       fh.Filename,
	}

	res, err := h.service.CreateBanner(ctx, description, uploadFile)
	if err != nil {
		log.Error().Err(err).Msg("handler::createBanner - Failed to create banner")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
