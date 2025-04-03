package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *productSubCategoryService) CreateProductSubCategory(ctx context.Context, req *dto.CreateOrUpdateProductSubCategoryRequest, categoryID int) (*dto.CreateOrUpdateProductSubCategoryResponse, error) {
	// check find product by category id
	productResult, err := s.productCategoryRepository.FindProductCategoryByID(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::CreateProductSubCategory - product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::CreateProductSubCategory - error finding product category by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert new product sub category
	res, err := s.productSubCategoryRepository.InsertNewProductSubCategory(ctx, req.SubCategory, productResult.ID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductSubCategoryAlreadyRegistered) {
			log.Error().Err(err).Msg("service::CreateProductSubCategory - product sub category already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrProductCategoryAlreadyRegistered))
		}

		log.Error().Err(err).Msg("service::CreateProductSubCategory - error inserting new product sub category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// create response
	response := dto.CreateOrUpdateProductSubCategoryResponse{
		ID:                res.ID,
		SubCategory:       res.Name,
		ProductCategoryID: res.ProductCategoryID,
	}

	// return response
	return &response, nil
}

func (s *productSubCategoryService) UpdateProductSubCategory(ctx context.Context, req *dto.CreateOrUpdateProductSubCategoryRequest, categoryID, subCategoryID int) (*dto.CreateOrUpdateProductSubCategoryResponse, error) {
	// check find product by category id
	_, err := s.productCategoryRepository.FindProductCategoryByID(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::UpdateProductSubCategory - product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateProductSubCategory - error finding product category by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check find sub product category by id
	productCategoryResult, err := s.productSubCategoryRepository.FindProductSubCategoryByID(ctx, subCategoryID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductSubCategoryNotFound) {
			log.Error().Err(err).Msg("service::UpdateProductSubCategory - product sub category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductSubCategoryNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateProductSubCategory - error finding product sub category by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// update product sub category
	res, err := s.productSubCategoryRepository.UpdateProducSubCategory(ctx, req.SubCategory, productCategoryResult.ID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductSubCategoryAlreadyRegistered) {
			log.Warn().Msg("service::UpdateProductSubCategory - product sub category already registered, returning existing record")
			response := dto.CreateOrUpdateProductSubCategoryResponse{
				ID:                productCategoryResult.ID,
				SubCategory:       productCategoryResult.Name,
				ProductCategoryID: productCategoryResult.ProductCategoryID,
			}

			return &response, nil
		}

		log.Error().Err(err).Msg("service::UpdateProductSubCategory - error updating product sub category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// create response
	response := dto.CreateOrUpdateProductSubCategoryResponse{
		ID:                res.ID,
		SubCategory:       res.Name,
		ProductCategoryID: res.ProductCategoryID,
	}

	// return response
	return &response, nil
}
