package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
)

func (s *productCategoryService) GetListProductCategory(ctx context.Context, page, limit int, search string) (*dto.GetListProductCategory, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list product category
	productCategories, total, err := s.productCategoryRepository.FindListProductCategory(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProductCategory - error getting list product category")
		return nil, err
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

func (s *productCategoryService) CreateProductCategory(ctx context.Context, name string) (*dto.CreateProductCategoryResponse, error) {
	// insert new product category
	res, err := s.productCategoryRepository.InsertNewProductCategory(ctx, name)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateProductCategory - error inserting new product category")
		return nil, err
	}

	// create response
	response := dto.CreateProductCategoryResponse{
		ID:       res.ID,
		Category: res.Name,
	}

	// return response
	return &response, nil
}
