package dto

import (
	productCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
)

type CreateOrUpdateProductSubCategoryRequest struct {
	SubCategory string `json:"subcategory" validate:"required,min=3,max=50,xss_safe"`
}

type CreateOrUpdateProductSubCategoryResponse struct {
	ID                int    `json:"id"`
	SubCategory       string `json:"subcategory"`
	ProductCategoryID int    `json:"category_id"`
}

type GetListProductSubCategory struct {
	SubCategory []GetListSubCategory `json:"subcategory"`
	TotalPages  int                  `json:"total_page"`
	CurrentPage int                  `json:"current_page"`
	PageSize    int                  `json:"page_size"`
	TotalData   int                  `json:"total_data"`
}

type GetListSubCategory struct {
	ID          int                             `json:"id"`
	SubCategory string                          `json:"subcategory"`
	Category    productCategory.GetListCategory `json:"category"`
}

type ProductSubCategory struct {
	ID          int    `json:"id"`
	SubCategory string `json:"subcategory"`
	CategoryID  int    `json:"category_id"`
}
