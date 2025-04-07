package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	InsertNewProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error)
	CountProductByName(ctx context.Context, name string) (int, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.ProductData, payloadFiles []dto.UploadFileRequest) (*dto.CreateOrUpdateProductResponse, error)
}
