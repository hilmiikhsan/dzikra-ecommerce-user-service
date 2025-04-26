package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
)

type SubDistrictService interface {
	GetListSubDistrict(ctx context.Context, districtID int) ([]dto.GetListSubDistrictResponse, error)
}
