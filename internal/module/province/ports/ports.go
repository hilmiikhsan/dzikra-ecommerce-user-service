package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
)

type ProvinceService interface {
	GetListProvince(ctx context.Context) ([]dto.GetListProvinceResponse, error)
}
