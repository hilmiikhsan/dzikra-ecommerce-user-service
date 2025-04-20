package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/entity"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	InsertNewProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, tx *sqlx.Tx, id int, data *entity.Product) (*entity.Product, error)
	CountProductByName(ctx context.Context, name string) (int, error)
	FindListProduct(ctx context.Context, limit, offset int, search string, categoryID, subcategoryID int) ([]dto.GetListProduct, int, error)
	FindProductByID(ctx context.Context, id int) (*entity.Product, error)
	SoftDeleteProductByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.ProductData, payloadFiles []dto.UploadFileRequest) (*dto.CreateOrUpdateProductResponse, error)
	UpdateProduct(ctx context.Context, productID int, req *dto.ProductData, payloadFiles []dto.UploadFileRequest) (*dto.CreateOrUpdateProductResponse, error)
	GetListProduct(ctx context.Context, page, limit int, search string, categoryID, subcategoryID int) (*dto.GetListProductResponse, error)
	GetDetailProduct(ctx context.Context, id int) (*dto.GetListProduct, error)
	RemoveProduct(ctx context.Context, id int) error
}
