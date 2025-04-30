package ports

import (
	"context"

	city "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
	province "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	subDistrict "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
)

type RajaongkirService interface {
	GetListCity(ctx context.Context, provinceID int) ([]city.GetListCityResponse, error)
	GetListProvince(ctx context.Context) ([]province.GetListProvinceResponse, error)
	GetListSubDistrict(ctx context.Context, districtID int) ([]subDistrict.GetListSubDistrictResponse, error)
}
