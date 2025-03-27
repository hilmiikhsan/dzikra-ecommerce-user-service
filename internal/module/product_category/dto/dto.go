package dto

type GetListProductCategory struct {
	Category    []GetListCategory `json:"category"`
	TotalPages  int               `json:"total_pages"`
	CurrentPage int               `json:"current_page"`
	PageSize    int               `json:"page_size"`
	TotalData   int               `json:"total_data"`
}

type GetListCategory struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
}

type CreateProductCategoryRequest struct {
	Category string `json:"category" validate:"required,min=3,max=50,xss_safe"`
}

type CreateProductCategoryResponse struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
}
