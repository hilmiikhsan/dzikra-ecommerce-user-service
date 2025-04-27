package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	InsertNewAddress(ctx context.Context, tx *sqlx.Tx, data *entity.Address) (*entity.Address, error)
	UpdateAddress(ctx context.Context, tx *sqlx.Tx, id int, data *entity.Address) (*entity.Address, error)
	SoftDeleteAddressByID(ctx context.Context, tx *sqlx.Tx, id int, userID uuid.UUID) error
}

type AddressService interface {
	CreateAddress(ctx context.Context, req *dto.CreateOrUpdateAddressRequest) (*dto.CreateOrUpdateAddressResponse, error)
	UpdateAddress(ctx context.Context, req *dto.CreateOrUpdateAddressRequest, addressID int) (*dto.CreateOrUpdateAddressResponse, error)
	RemoveAddress(ctx context.Context, addressID int, userID string) error
}
