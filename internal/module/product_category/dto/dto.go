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
