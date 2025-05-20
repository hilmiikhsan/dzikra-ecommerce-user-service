package dto

type CreateNotificationRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=100,xss_safe"`
	Detail  string `json:"detail" validate:"required,min=3,max=100,xss_safe"`
	Url     string `json:"url" validate:"required,url,max=100,xss_safe"`
	NTypeID string `json:"n_type_id" validate:"required,min=3,max=50,xss_safe"`
	UserID  string `json:"user_id" validate:"required,uuid,min=3,max=50,xss_safe"`
}

type GetListNotificationResponse struct {
	Notification []NotificationDetail `json:"notification"`
	TotalPages   int                  `json:"total_page"`
	CurrentPage  int                  `json:"current_page"`
	PageSize     int                  `json:"page_size"`
	TotalData    int                  `json:"total_data"`
}

type NotificationDetail struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Detail    *string `json:"detail"`
	Url       string  `json:"url"`
	NTypeID   string  `json:"n_type_id"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
}

type SendFcmBatchRequest struct {
	IsUser bool   `json:"isUser"`
	Title  string `json:"title" validate:"required,min=3,max=100,xss_safe"`
	Detail string `json:"detail" validate:"required,min=3,max=100,xss_safe"`
	Url    string `json:"url" validate:"required,url,max=100,xss_safe"`
}
