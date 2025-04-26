package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
)

type CityService interface {
	GetListCity(ctx context.Context, provinceID int) ([]dto.GetListCityResponse, error)
}
