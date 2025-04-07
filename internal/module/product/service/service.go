package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	product "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	productGroceryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	productImageDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	productVariantDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/dto"
	productVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *productService) CreateProduct(ctx context.Context, req *dto.ProductData, payloadFiles []dto.UploadFileRequest) (*dto.CreateOrUpdateProductResponse, error) {
	// convert categoryID to int
	productCategoryID, err := strconv.Atoi(req.CategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - error converting categoryID to int")
		return nil, err
	}

	// check product category if exist
	productCategoryCount, err := s.productCategoryRepository.CountProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - error checking product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if productCategoryCount == 0 {
		log.Error().Err(err).Msg("service::CreateProduct - product category not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
	}

	// convert subCategoryID to int
	subCategoryID, err := strconv.Atoi(req.SubCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - error converting subCategoryID to int")
		return nil, err
	}

	// check product sub category if exist
	subCategoryCount, err := s.productSubCategoryRepository.CountProductSubCategoryByID(ctx, subCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - error checking product sub category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if subCategoryCount == 0 {
		log.Error().Err(err).Msg("service::CreateProduct - product sub category not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductSubCategoryNotFound))
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateProduct - Failed to rollback transaction")
			}
		}
	}()

	// insert new product
	productResult, err := s.productRepository.InsertNewProduct(ctx, tx, &product.Product{
		Name:                 req.Name,
		RealPrice:            req.RealPrice,
		DiscountPrice:        req.DiscountPrice,
		CapitalPrice:         req.CapitalPrice,
		Description:          req.Description,
		Spesification:        req.Spec,
		Stock:                req.Stock,
		Weight:               req.Weight,
		VariantName:          req.VariantName,
		ProductCategoryID:    productCategoryID,
		ProductSubCategoryID: subCategoryID,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - error inserting new product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert new product variant
	var productVariantResults []*productVariant.ProductVariant
	for _, variant := range req.Variants {
		variantResult, err := s.productVariantRepository.InsertNewProductVariant(ctx, tx, &productVariant.ProductVariant{
			VariantSubName: variant.VariantSubName,
			VariantStock:   variant.VariantStock,
			VariantWeight:  variant.VariantWeight,
			CapitalPrice:   variant.CapitalPrice,
			RealPrice:      variant.RealPrice,
			DiscountPrice:  variant.DiscountPrice,
			ProductID:      productResult.ID,
		})
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - error inserting new product variant")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		productVariantResults = append(productVariantResults, variantResult)
	}

	// insert new product grocery
	var productGroceryResults []*productGrocery.ProductGrocery
	for _, grocery := range req.GroceryPrices {
		groceryResult, err := s.productGroceryRepository.InsertNewProductGrocery(ctx, tx, &productGrocery.ProductGrocery{
			MinBuy:    grocery.MinBuy,
			Discount:  grocery.Discount,
			ProductID: productResult.ID,
		})
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateProduct - error inserting new product grocery")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		productGroceryResults = append(productGroceryResults, groceryResult)
	}

	// check current total image
	currentImageCount, err := s.productImageRepository.CountProductImagesByProductID(ctx, productResult.ID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - error counting product images")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert new product image
	var insertedImages []productImage.ProductImage
	for i := range payloadFiles {
		sortOrder := currentImageCount + i + 1

		result, err := s.productImageRepository.InsertNewProductImage(ctx, tx, &productImage.ProductImage{
			ImageURL:  "",
			ProductID: productResult.ID,
			Sort:      sortOrder,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::CreateProduct - error inserting product image record")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		insertedImages = append(insertedImages, *result)
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateProduct - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload files to minio
	for i, fileReq := range payloadFiles {
		ext := strings.ToLower(filepath.Ext(fileReq.Filename))
		sortOrder := insertedImages[i].Sort
		objectName := fmt.Sprintf("product_images/%s_%d%s", utils.GenerateBucketFileUUID(), sortOrder, ext)
		byteFile := utils.NewByteFile(fileReq.File)

		uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, fileReq.FileHeaderSize, fileReq.ContentType)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateProduct - Failed to upload file to minio")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		updatedImage, err := s.productImageRepository.UpdateProductImageURL(ctx, insertedImages[i].ID, objectName)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateProduct - Failed to update product image URL")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		insertedImages[i] = *updatedImage

		log.Info().Msgf("Uploaded image URL: %s", uploadedPath)
	}

	// create response
	response := &dto.CreateOrUpdateProductResponse{
		ID:             productResult.ID,
		Name:           req.Name,
		Description:    req.Description,
		Spec:           req.Spec,
		RealPrice:      req.RealPrice,
		CapitalPrice:   req.CapitalPrice,
		DiscountPrice:  req.DiscountPrice,
		Stock:          req.Stock,
		Weight:         req.Weight,
		CategoryID:     productCategoryID,
		SubCategoryID:  subCategoryID,
		VariantName:    req.VariantName,
		ProductVariant: nil,
		ProductGrocery: nil,
		ProductImage:   nil,
	}

	// Map product variants to response
	for _, variantResult := range productVariantResults {
		response.ProductVariant = append(response.ProductVariant, productVariantDto.ProductVariant{
			ID:             variantResult.ID,
			VariantSubName: variantResult.VariantSubName,
			VariantStock:   variantResult.VariantStock,
			VariantWeight:  variantResult.VariantWeight,
			CapitalPrice:   variantResult.CapitalPrice,
			RealPrice:      variantResult.RealPrice,
			DiscountPrice:  variantResult.DiscountPrice,
			ProductID:      productResult.ID,
		})
	}

	// Map product groceries to response
	for _, groceryResult := range productGroceryResults {
		response.ProductGrocery = append(response.ProductGrocery, productGroceryDto.ProductGrocery{
			ID:        groceryResult.ID,
			MinBuy:    groceryResult.MinBuy,
			Discount:  groceryResult.Discount,
			ProductID: productResult.ID,
		})
	}

	// Map product images to response
	publicURL := config.Envs.MinioStorage.PublicURL
	for _, image := range insertedImages {
		response.ProductImage = append(response.ProductImage, productImageDto.ProductImage{
			ID:        image.ID,
			ImageURL:  utils.FormatMediaPathURL(image.ImageURL, publicURL),
			Position:  image.Sort,
			ProductID: image.ProductID,
		})
	}

	return response, nil
}
