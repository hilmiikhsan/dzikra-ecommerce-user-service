package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *cartService) AddCartItem(ctx context.Context, req *dto.AddCartItemRequest) (*dto.AddCartItemResponse, error) {
	// check count product data
	productCount, err := s.productRepository.CountProductByID(ctx, req.ProductID)
	if err != nil {
		log.Error().Err(err).Msg("service::AddCartItem - failed to count product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if productCount == 0 {
		log.Warn().Msg("service::AddCartItem - product not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
	}

	// check count product variant data
	productVariantCount, err := s.productVariantRepository.CountProductVariantByIDAndProductID(ctx, req.ProductVariantID, req.ProductID)
	if err != nil {
		log.Error().Err(err).Msg("service::AddCartItem - failed to count product variant")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if productVariantCount == 0 {
		log.Warn().Msg("service::AddCartItem - product variant not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductVariantsNotFound))
	}

	// convert user id to int
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("service::AddCartItem - failed to parse user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::AddCartItem - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::AddCartItem - Failed to rollback transaction")
			}
		}
	}()

	// insert new cart
	res, err := s.cartRepository.InsertNewCart(ctx, tx, &entity.Cart{
		UserID:           userID,
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		Quantity:         req.Quantity,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::AddCartItem - Failed to insert new cart")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::AddCartItem - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.AddCartItemResponse{
		ID:               res.ID,
		UserID:           res.UserID.String(),
		ProductID:        res.ProductID,
		ProductVariantID: res.ProductVariantID,
		Quantity:         res.Quantity,
		CreatedAt:        utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *cartService) GetListCart(ctx context.Context, userID string) (*[]dto.GetListCartResponse, error) {
	// convert user id to uuid
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListCart - failed to parse user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get cart by user id
	carts, err := s.cartRepository.FindListCartByUserID(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListCart - failed to get cart by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if carts is nil
	if carts == nil {
		carts = []dto.GetListCartResponse{}
	}

	return &carts, nil
}
