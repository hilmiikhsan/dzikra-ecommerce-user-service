package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/ports"
	productGroceryDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productImageDto "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CartRepository = &cartRepository{}

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *cartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) InsertNewCart(ctx context.Context, tx *sqlx.Tx, data *entity.Cart) (*entity.Cart, error) {
	var res = new(entity.Cart)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewCart),
		data.UserID,
		data.ProductID,
		data.ProductVariantID,
		data.Quantity,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.ProductID,
		&res.ProductVariantID,
		&res.Quantity,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewCart - Failed to insert new cart")
		return nil, err
	}

	return res, nil
}

func (r *cartRepository) FindListCartByUserID(ctx context.Context, userID uuid.UUID) ([]dto.GetListCartResponse, error) {
	var rows []entity.Cart
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindListCartByUserID), userID); err != nil {
		log.Error().Err(err).Msg("repository::FindListCartByUserID - error executing query")
		return nil, fmt.Errorf("error querying cart: %w", err)
	}

	cartMap := make(map[int]*dto.GetListCartResponse, len(rows))
	ordered := make([]*dto.GetListCartResponse, 0, len(rows))

	for _, row := range rows {
		crt, exists := cartMap[row.ID]
		if !exists {
			crt = &dto.GetListCartResponse{
				ID:                          row.ID,
				Quantity:                    row.Quantity,
				ProductID:                   row.ProductID,
				ProductVariantID:            row.ProductVariantID,
				ProductName:                 row.ProductName,
				ProductRealPrice:            row.ProductRealPrice,
				ProductDiscountPrice:        row.ProductDiscountPrice,
				ProductStock:                row.ProductStock,
				ProductVariantName:          row.ProductVariantName,
				ProductGrocery:              []productGroceryDto.ProductGrocery{},
				ProductVariantSubName:       row.ProductVariantSubName.String,
				ProductVariantRealPrice:     row.ProductVariantRealPrice.String,
				ProductVariantDiscountPrice: row.ProductVariantDiscountPrice.String,
				ProductVariantStock:         int(row.ProductVariantStock.Int64),
				ProductImage:                []productImageDto.ProductImage{},
			}

			cartMap[row.ID] = crt
			ordered = append(ordered, crt)
		}

		if row.ProductGroceryID.Valid {
			gid := int(row.ProductGroceryID.Int64)
			dup := false
			for _, g := range crt.ProductGrocery {
				if g.ID == gid {
					dup = true
					break
				}
			}
			if !dup {
				crt.ProductGrocery = append(crt.ProductGrocery, productGroceryDto.ProductGrocery{
					ID:        gid,
					MinBuy:    int(row.ProductGroceryMinBuy.Int64),
					Discount:  int(row.ProductGroceryDiscount.Int64),
					ProductID: int(row.ProductGroceryProductID.Int64),
				})
			}
		}

		if row.ProductImageID.Valid {
			iid := int(row.ProductImageID.Int64)
			dup := false
			for _, img := range crt.ProductImage {
				if img.ID == iid {
					dup = true
					break
				}
			}
			if !dup {
				crt.ProductImage = append(crt.ProductImage, productImageDto.ProductImage{
					ID:        iid,
					ImageURL:  row.ProductImageURL.String,
					Position:  int(row.ProductImageSort.Int64),
					ProductID: int(row.ProductImageProductID.Int64),
				})
			}
		}
	}

	result := make([]dto.GetListCartResponse, len(ordered))
	for i, p := range ordered {
		result[i] = *p
	}

	return result, nil
}

func (r *cartRepository) UpdateCart(ctx context.Context, tx *sqlx.Tx, data *entity.Cart) (*entity.Cart, error) {
	var res = new(entity.Cart)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateCart),
		data.UserID,
		data.ProductID,
		data.ProductVariantID,
		data.Quantity,
		data.ID,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.ProductID,
		&res.ProductVariantID,
		&res.Quantity,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateCart - cart with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrCartNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateCart - Failed to update cart")
		return nil, err
	}

	return res, nil
}

func (r *cartRepository) DeleteCartByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := r.db.ExecContext(ctx, r.db.Rebind(queryDeleteCartByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::DeleteCartByID - Failed to soft delete cart")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::DeleteCartByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrCartNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::DeleteCartByID - Cart not found")
		return errNotFound
	}

	return nil
}
