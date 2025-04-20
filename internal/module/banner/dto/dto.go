package dto

type UploadFileRequest struct {
	ObjectName     string `json:"object_name"`
	File           []byte `json:"-"`
	FileHeaderSize int64  `json:"-"`
	ContentType    string `json:"-"`
	Filename       string `json:"-"`
}

type CreateBannerResponse struct {
	ID          int    `json:"id"`
	ImageURL    string `json:"image_url"`
	Description string `json:"desc"`
}

type GetListBannerResponse struct {
	Banner      []GetListBanner `json:"banner"`
	TotalPages  int             `json:"total_page"`
	CurrentPage int             `json:"current_page"`
	PageSize    int             `json:"page_size"`
	TotalData   int             `json:"total_data"`
}

type GetListBanner struct {
	ID          int    `json:"id"`
	ImageURL    string `json:"image_url"`
	Description string `json:"desc"`
}
