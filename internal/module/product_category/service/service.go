package service

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
)

func (s *productCategoryService) GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategory, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list product category
	productCategories, total, err := s.productCategoryRepository.FindListProductCategory(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProductCategory - error getting list product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if productCategories is nil
	if productCategories == nil {
		productCategories = []dto.GetListCategory{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListProductCategory{
		Category:    productCategories,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *productCategoryService) CreateProductCategory(ctx context.Context, name string) (*dto.CreateOrProductCategoryResponse, error) {
	// insert new product category
	res, err := s.productCategoryRepository.InsertNewProductCategory(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryAlreadyRegistered) {
			log.Error().Err(err).Msg("service::CreateProductCategory - product category already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrProductCategoryAlreadyRegistered))
		}

		log.Error().Err(err).Msg("service::CreateProductCategory - error inserting new product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// create response
	response := dto.CreateOrProductCategoryResponse{
		ID:       res.ID,
		Category: res.Name,
	}

	// return response
	return &response, nil
}

func (s *productCategoryService) UpdateProductCategory(ctx context.Context, id int, name string) (*dto.CreateOrProductCategoryResponse, error) {
	// First, fetch the existing record.
	productCategoryResult, err := s.productCategoryRepository.FindProductCategoryByID(ctx, id)
	if err != nil {
		// If not found, return a Not Found error.
		if err == sql.ErrNoRows || strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::UpdateProductCategory - Product category not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}
		log.Error().Err(err).Msg("service::UpdateProductCategory - Error finding product category by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Attempt to update the product category.
	updatedCategory, err := s.productCategoryRepository.UpdateProductCategory(ctx, id, name)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductCategoryAlreadyRegistered) {
			log.Warn().Msg("service::UpdateProductCategory - product category already registered, returning existing record")
			response := dto.CreateOrProductCategoryResponse{
				ID:       productCategoryResult.ID,
				Category: productCategoryResult.Name,
			}
			return &response, nil
		}

		log.Error().Err(err).Msg("service::UpdateProductCategory - error updating product category")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Otherwise, update was successful; build the response using the updated record.
	response := dto.CreateOrProductCategoryResponse{
		ID:       updatedCategory.ID,
		Category: updatedCategory.Name,
	}

	// Return the response.
	return &response, nil
}

func (s *productCategoryService) RemoveProductCategory(ctx context.Context, id int) error {
	// Check if the product category exists.
	productCategory, err := s.productCategoryRepository.FindProductCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), constants.ErrProductCategoryNotFound) {
			log.Error().Err(err).Msg("service::RemoveProductCategory - Product category not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductCategoryNotFound))
		}
		log.Error().Err(err).Msg("service::RemoveProductCategory - Error finding product category by id")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Attempt to delete the product category.
	err = s.productCategoryRepository.DeleteProductCategoryByID(ctx, productCategory.ID)
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveProductCategory - error deleting product category")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
