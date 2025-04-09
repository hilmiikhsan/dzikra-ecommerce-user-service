package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
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

func (h *productHandler) updateProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	productIDStr := c.Params("product_id")
	if productIDStr == "" {
		log.Error().Msg("handler::updateProduct - product ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product ID is required"))
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateProduct - invalid product ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product ID"))
	}

	productDataStr := c.FormValue("product_data")
	if productDataStr == "" {
		log.Error().Msg("handler::updateProduct - product_data is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Product data is required"))
	}

	var productData dto.ProductData
	if err := json.Unmarshal([]byte(productDataStr), &productData); err != nil {
		log.Error().Err(err).Msg("handler::updateProduct - error unmarshalling product_data")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product data"))
	}

	if err := h.validator.Validate(productData); err != nil {
		log.Warn().Err(err).Msg("handler::updateProduct - invalid product data")
		code, errs := err_msg.Errors(err, &productData)
		return c.Status(code).JSON(response.Error(errs))
	}

	var uploadFiles []dto.UploadFileRequest
	mf, err := c.MultipartForm()
	if err == nil {
		fileHeaders := mf.File[constants.MultipartFormFile]
		for _, fh := range fileHeaders {
			if fh.Size > constants.MaxFileSize {
				log.Error().Msg("handler::updateProduct - file size exceeds limit")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("File size exceeds limit"))
			}

			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !constants.AllowedImageExtensions[ext] {
				log.Error().Msg("handler::updateProduct - invalid file extension")
				return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid file extension"))
			}

			file, err := fh.Open()
			if err != nil {
				log.Error().Err(err).Msg("handler::updateProduct - failed to open file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}

			fileBytes, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				log.Error().Err(err).Msg("handler::updateProduct - failed to read file")
				return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error"))
			}

			mimeType := http.DetectContentType(fileBytes)
			if !strings.HasPrefix(mimeType, "image/") {
				log.Error().Msg("handler::updateProduct - uploaded file is not a valid image")
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
	}

	if len(uploadFiles) == 0 && len(productData.ImageToKeep) == 0 {
		log.Error().Msg("handler::updateProduct - No image file provided (neither new upload nor images to keep)")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("At least one image file must be provided"))
	}

	res, err := h.service.UpdateProduct(ctx, productID, &productData, uploadFiles)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateProduct - failed to update product")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productHandler) getListProduct(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListProduct(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProduct - failed to get list of products")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
