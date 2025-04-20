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
