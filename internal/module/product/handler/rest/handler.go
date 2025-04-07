package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productHandler) createProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	productDataStr := c.FormValue("product_data")
	if productDataStr == "" {
		log.Error().Msg("handler::createProduct - product_data is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product data is required"))
	}

	var productData dto.ProductData
	if err := json.Unmarshal([]byte(productDataStr), &productData); err != nil {
		log.Error().Err(err).Msg("handler::createProduct - error unmarshalling product_data")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product data"))
	}

	if err := h.validator.Validate(productData); err != nil {
		log.Warn().Err(err).Msg("handler::createProduct - Invalid product data")
		code, errs := err_msg.Errors(err, &productData)
		return c.Status(code).JSON(response.Error(errs))
	}

	mf, err := c.MultipartForm()
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to parse multipart form")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid multipart form data"))
	}

	fileHeaders := mf.File[constants.MultipartFormFile]
	if len(fileHeaders) == 0 {
		log.Error().Msg("handler::createProduct - No image files uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file is required"))
	}

	var uploadFiles []dto.UploadFileRequest
	for _, fh := range fileHeaders {
		if fh.Size > constants.MaxFileSize {
			log.Error().Msg("handler::createProduct - File size exceeds limit")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
		}

		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if !constants.AllowedImageExtensions[ext] {
			log.Error().Msg("handler::createProduct - Invalid file extension")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
		}

		file, err := fh.Open()
		if err != nil {
			log.Error().Err(err).Msg("handler::createProduct - Failed to open file")
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
		}

		fileBytes, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			log.Error().Err(err).Msg("handler::createProduct - Failed to read file")
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
		}

		mimeType := http.DetectContentType(fileBytes)
		if !strings.HasPrefix(mimeType, "image/") {
			log.Error().Msg("handler::createProduct - Uploaded file is not a valid image")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Uploaded file is not a valid image"))
		}

		objectName := "product_images/" + utils.GenerateBucketFileUUID() + ext

		uploadFiles = append(uploadFiles, dto.UploadFileRequest{
			ObjectName:     objectName,
			File:           fileBytes,
			FileHeaderSize: fh.Size,
			ContentType:    mimeType,
			Filename:       fh.Filename,
		})
	}

	res, err := h.service.CreateProduct(ctx, &productData, uploadFiles)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProduct - Failed to create product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
