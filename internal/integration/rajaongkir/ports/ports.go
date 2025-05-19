package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/dto"
	address "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	city "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
	province "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	subDistrict "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
)

type RajaongkirService interface {
	GetListCity(ctx context.Context, provinceID int) ([]city.GetListCityResponse, error)
	GetListProvince(ctx context.Context) ([]province.GetListProvinceResponse, error)
	GetListSubDistrict(ctx context.Context, districtID int) ([]subDistrict.GetListSubDistrictResponse, error)
	GetShippingCost(ctx context.Context, weight, courier string, address *address.Address) ([]dto.CostResult, error)
	GetWaybill(ctx context.Context, waybill, courier string) (*dto.GetWaybillResponse, error)
}
