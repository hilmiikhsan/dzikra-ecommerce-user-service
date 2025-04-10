package service

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	product "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	productCategoryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	productGroceryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	productImageDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	productSubCategoryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	productVariantDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/dto"
	productVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
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
		Specification:  req.Spec,
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

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateProductResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productService) UpdateProduct(ctx context.Context, productID int, req *dto.ProductData, payloadFiles []dto.UploadFileRequest) (*dto.CreateOrUpdateProductResponse, error) {
	// convert categoryID to int
	productCategoryID, err := strconv.Atoi(req.CategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error converting categoryID to int")
		return nil, err
	}

	// check product category if exist
	productCategoryCount, err := s.productCategoryRepository.CountProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error checking product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if productCategoryCount == 0 {
		log.Error().Msg("service::UpdateProduct - product category not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
	}

	// convert subCategoryID ke integer
	subCategoryID, err := strconv.Atoi(req.SubCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error converting subCategoryID to int")
		return nil, err
	}

	// check product sub category if exist
	subCategoryCount, err := s.productSubCategoryRepository.CountProductSubCategoryByID(ctx, subCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error checking product sub category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if subCategoryCount == 0 {
		log.Error().Msg("service::UpdateProduct - product sub category not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductSubCategoryNotFound))
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("service::UpdateProduct - Failed to rollback transaction")
			}
		}
	}()

	// update product
	updatedProduct, err := s.productRepository.UpdateProduct(ctx, tx, productID, &product.Product{
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
		log.Error().Err(err).Msg("service::UpdateProduct - error updating product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// update or insert product variant
	var productVariantResults []*productVariant.ProductVariant
	for _, variant := range req.Variants {
		if variant.ID > 0 {
			// update existing product variant
			updatedVariant, err := s.productVariantRepository.UpdateProductVariant(ctx, tx, &productVariant.ProductVariant{
				ID:             variant.ID,
				VariantSubName: variant.VariantSubName,
				VariantStock:   variant.VariantStock,
				VariantWeight:  variant.VariantWeight,
				CapitalPrice:   variant.CapitalPrice,
				RealPrice:      variant.RealPrice,
				DiscountPrice:  variant.DiscountPrice,
				ProductID:      productID,
			})
			if err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - error updating product variant")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
			productVariantResults = append(productVariantResults, updatedVariant)
		} else {
			// insert new product variant
			newVariant, err := s.productVariantRepository.InsertNewProductVariant(ctx, tx, &productVariant.ProductVariant{
				VariantSubName: variant.VariantSubName,
				VariantStock:   variant.VariantStock,
				VariantWeight:  variant.VariantWeight,
				CapitalPrice:   variant.CapitalPrice,
				RealPrice:      variant.RealPrice,
				DiscountPrice:  variant.DiscountPrice,
				ProductID:      productID,
			})
			if err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - error inserting new product variant")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
			productVariantResults = append(productVariantResults, newVariant)
		}
	}

	// delete product variant if variant.ID is in req.DelVariants
	for _, delID := range req.DelVariants {
		if err := s.productVariantRepository.DeleteProductVariant(ctx, tx, delID, productID); err != nil {
			log.Error().Err(err).Msg("service::UpdateProduct - error deleting product variant")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// update or insert product grocery
	var productGroceryResults []*productGrocery.ProductGrocery
	for _, grocery := range req.GroceryPrices {
		if grocery.ID > 0 {
			updatedGrocery, err := s.productGroceryRepository.UpdateProductGrocery(ctx, tx, &productGrocery.ProductGrocery{
				ID:        grocery.ID,
				MinBuy:    grocery.MinBuy,
				Discount:  grocery.Discount,
				ProductID: productID,
			})
			if err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - error updating product grocery")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
			productGroceryResults = append(productGroceryResults, updatedGrocery)
		} else {
			newGrocery, err := s.productGroceryRepository.InsertNewProductGrocery(ctx, tx, &productGrocery.ProductGrocery{
				MinBuy:    grocery.MinBuy,
				Discount:  grocery.Discount,
				ProductID: productID,
			})
			if err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - error inserting new product grocery")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
			productGroceryResults = append(productGroceryResults, newGrocery)
		}
	}

	// delete product grocery if grocery.ID is in req.DelGrocery
	existingImages, err := s.productImageRepository.FindProductImagesByProductID(ctx, productID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error fetching product images")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	for _, img := range existingImages {
		keep := false
		for _, idToKeep := range req.ImageToKeep {
			if img.ID == idToKeep {
				keep = true
				break
			}
		}
		if !keep {
			if err := s.productImageRepository.DeleteProductImage(ctx, tx, img.ID); err != nil {
				log.Error().Err(err).Msg("service::UpdateProduct - error deleting product image")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}
	}

	// Insert new product images
	var newImageRecords []productImage.ProductImage
	for i := range payloadFiles {
		sortOrder := len(existingImages) + i + 1
		imageRecord, err := s.productImageRepository.InsertNewProductImage(ctx, tx, &productImage.ProductImage{
			ImageURL:  "",
			ProductID: productID,
			Sort:      sortOrder,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateProduct - error inserting product image record")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
		newImageRecords = append(newImageRecords, *imageRecord)
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload files to minio
	for i, fileReq := range payloadFiles {
		ext := strings.ToLower(filepath.Ext(fileReq.Filename))
		objectName := fmt.Sprintf("product_images/%s_%d%s", utils.GenerateBucketFileUUID(), newImageRecords[i].Sort, ext)
		byteFile := utils.NewByteFile(fileReq.File) // ByteFile mengembalikan objek io.ReadSeeker

		_, err := s.minioService.UploadFile(ctx, objectName, byteFile, fileReq.FileHeaderSize, fileReq.ContentType)
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateProduct - Failed to upload file to minio")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		updatedImage, err := s.productImageRepository.UpdateProductImageURL(ctx, newImageRecords[i].ID, objectName)
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateProduct - Failed to update product image URL")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
		newImageRecords[i] = *updatedImage

		log.Info().Msgf("Uploaded image URL: %s", objectName)
	}

	// Mapping response
	response := &dto.CreateOrUpdateProductResponse{
		ID:             updatedProduct.ID,
		Name:           req.Name,
		Description:    req.Description,
		Specification:  req.Spec,
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
			ProductID:      productID,
		})
	}

	// Map product groceries to response
	for _, groceryResult := range productGroceryResults {
		response.ProductGrocery = append(response.ProductGrocery, productGroceryDto.ProductGrocery{
			ID:        groceryResult.ID,
			MinBuy:    groceryResult.MinBuy,
			Discount:  groceryResult.Discount,
			ProductID: productID,
		})
	}

	// Map product images to response
	publicURL := config.Envs.MinioStorage.PublicURL
	var finalProductImages []productImageDto.ProductImage
	existingImages, err = s.productImageRepository.FindProductImagesByProductID(ctx, productID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateProduct - error fetching product images for response")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	for _, image := range existingImages {
		finalProductImages = append(finalProductImages, productImageDto.ProductImage{
			ID:        image.ID,
			ImageURL:  utils.FormatMediaPathURL(image.ImageURL, publicURL),
			Position:  image.Sort,
			ProductID: image.ProductID,
		})
	}
	response.ProductImage = finalProductImages

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateProductResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productService) GetListProduct(ctx context.Context, page, limit int, search string) (*dto.GetListProductResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list product
	products, total, err := s.productRepository.FindListProduct(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProduct - error getting list product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if products is nil
	if products == nil {
		products = []dto.GetListProduct{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListProductResponse{
		Product:     products,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *productService) GetDetailProduct(ctx context.Context, id int) (*dto.GetListProduct, error) {
	// find product by id
	productResult, err := s.productRepository.FindProductByID(ctx, id)
	if err != nil {
		// check if product not found
		if err == sql.ErrNoRows || strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::GetDetailProduct - Product not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailProduct - Error finding product by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Map entity to DTO
	response := &dto.GetListProduct{
		ID:            productResult.ID,
		Name:          productResult.Name,
		Description:   productResult.Description,
		Specification: productResult.Spesification,
		RealPrice:     productResult.RealPrice,
		CapitalPrice:  productResult.CapitalPrice,
		DiscountPrice: productResult.DiscountPrice,
		Stock:         productResult.Stock,
		Weight:        productResult.Weight,
		VariantName:   productResult.VariantName,
		ProductCategory: productCategoryDto.GetListCategory{
			ID:       productResult.ProductCategoryID,
			Category: "",
		},
		ProductSubCategory: productSubCategoryDto.ProductSubCategory{
			ID:          productResult.ProductSubCategoryID,
			SubCategory: "", // Sama seperti di atas
			CategoryID:  productResult.ProductSubCategoryID,
		},
		ProductVariant: []productVariantDto.ProductVariant{},
		ProductGrocery: []productGroceryDto.ProductGrocery{},
		ProductImage:   []productImageDto.ProductImage{},
	}

	// Map variant
	for _, variant := range productResult.ProductVariant {
		response.ProductVariant = append(response.ProductVariant, productVariantDto.ProductVariant{
			ID:             variant.ID,
			VariantSubName: variant.VariantSubName,
			VariantStock:   variant.VariantStock,
			VariantWeight:  variant.VariantWeight,
			CapitalPrice:   variant.CapitalPrice,
			RealPrice:      variant.RealPrice,
			DiscountPrice:  variant.DiscountPrice,
			ProductID:      variant.ProductID,
		})
	}

	// Map groceries
	for _, grocery := range productResult.ProductGrocery {
		response.ProductGrocery = append(response.ProductGrocery, productGroceryDto.ProductGrocery{
			ID:        grocery.ID,
			MinBuy:    grocery.MinBuy,
			Discount:  grocery.Discount,
			ProductID: grocery.ProductID,
		})
	}

	// Map images
	publicURL := config.Envs.MinioStorage.PublicURL
	for _, image := range productResult.ProductImage {
		response.ProductImage = append(response.ProductImage, productImageDto.ProductImage{
			ID:        image.ID,
			ImageURL:  utils.FormatMediaPathURL(image.ImageURL, publicURL),
			Position:  image.Sort,
			ProductID: image.ProductID,
		})
	}

	// Sanitasi response untuk menghindari XSS
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeGetListProduct(*response, policy)

	return &sanitizedResponse, nil
}

func (s *productService) RemoveProduct(ctx context.Context, id int) error {
	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("service::RemoveProduct - Failed to rollback transaction")
			}
		}
	}()

	// soft delete product
	if err := s.productRepository.SoftDeleteProductByID(ctx, tx, id); err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Msg("service::RemoveProduct - Product not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveProduct - error soft deleting product")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete product variants
	if err := s.productVariantRepository.SoftDeleteProductVariantsByProductID(ctx, tx, id); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - error soft deleting product variants")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete product groceries
	if err := s.productGroceryRepository.SoftDeleteProductGroceriesByProductID(ctx, tx, id); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - error soft deleting product groceries")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if product images exist
	existingImages, err := s.productImageRepository.FindProductImagesByProductID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("service::DeleteProduct - error fetching product images")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete product images
	if err := s.productImageRepository.SoftDeleteProductImagesByProductID(ctx, tx, id); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - error soft deleting product images")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveProduct - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// delete files from minio
	for _, img := range existingImages {
		if err := s.minioService.DeleteFile(ctx, img.ImageURL); err != nil {
			log.Error().Err(err).Msgf("service::DeleteProduct - failed to delete file '%s' from minio", img.ImageURL)
		} else {
			log.Info().Msgf("Successfully deleted file '%s' from minio", img.ImageURL)
		}
	}

	return nil
}
