package dto

type CreateOrUpdateProductSubCategoryRequest struct {
	SubCategory string `json:"subcategory" validate:"required,min=3,max=50,xss_safe"`
}

type CreateOrUpdateProductSubCategoryResponse struct {
	ID                int    `json:"id"`
	SubCategory       string `json:"subcategory"`
	ProductCategoryID int    `json:"category_id"`
}
