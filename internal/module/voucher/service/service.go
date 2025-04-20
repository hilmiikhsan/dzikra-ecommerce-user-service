package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *voucherService) CreateVoucher(ctx context.Context, req *dto.CreateOrUpdateVoucherRequest) (*dto.CreateOrUpdateVoucherResponse, error) {
	// check if voucher type already exists
	countVoucherTypeResult, err := s.voucherTypeRepository.CountVoucherTypeByType(ctx, req.VoucherType)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrVoucherTypeNotFound) {
			log.Error().Err(err).Msg("service::CreateVoucher - voucher type not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrVoucherTypeNotFound))
		}

		log.Error().Err(err).Msg("service::CreateVoucher - error checking voucher type")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check startAt and endAt
	startAt, err := utils.ParseTime(req.StartAt)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateVoucher - error parsing startAt")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidStartAt))
	}
	endAt, err := utils.ParseTime(req.EndAt)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateVoucher - error parsing endAt")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidEndAt))
	}

	// check if voucher code already exists
	countVoucherCode, err := s.voucherRepository.InsertNewVoucher(ctx, &entity.Voucher{
		Name:          req.Name,
		VoucherQuota:  req.VoucherQuota,
		Code:          req.Code,
		Discount:      req.Discount,
		StartAt:       startAt,
		EndAt:         endAt,
		VoucherTypeID: countVoucherTypeResult.ID,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrVoucherCodeAlreadyRegistered) {
			log.Error().Err(err).Msg("service::CreateVoucher - voucher code already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrVoucherCodeAlreadyRegistered))
		}

		log.Error().Err(err).Msg("service::CreateVoucher - error inserting new voucher")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// create response
	res := &dto.CreateOrUpdateVoucherResponse{
		Name:          countVoucherCode.Name,
		VoucherQuota:  countVoucherCode.VoucherQuota,
		Code:          countVoucherCode.Code,
		Discount:      countVoucherCode.Discount,
		VoucherTypeID: countVoucherTypeResult.Type,
		CreatedAt:     utils.FormatTime(countVoucherCode.CreatedAt),
		StartAt:       utils.FormatToWIB(countVoucherCode.StartAt),
		EndAt:         utils.FormatToWIB(countVoucherCode.EndAt),
	}

	return res, nil
}

func (s *voucherService) GetListVoucher(ctx context.Context, page, limit int, search string) (*dto.GetListVoucherResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list voucher
	vouchers, total, err := s.voucherRepository.FindListVoucher(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListVoucher - error getting list voucher")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if vouchers is nil
	if vouchers == nil {
		vouchers = []dto.GetListVoucher{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListVoucherResponse{
		Voucher:     vouchers,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *voucherService) UpdateVoucher(ctx context.Context, id int, req *dto.CreateOrUpdateVoucherRequest) (*dto.CreateOrUpdateVoucherResponse, error) {
	// check if voucher type already exists
	countVoucherTypeResult, err := s.voucherTypeRepository.CountVoucherTypeByType(ctx, req.VoucherType)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrVoucherTypeNotFound) {
			log.Error().Err(err).Msg("service::UpdateVoucher - voucher type not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrVoucherTypeNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateVoucher - checking voucher type")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check startAt and endAt
	startAt, err := utils.ParseTime(req.StartAt)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateVoucher - error parsing startAt")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidStartAt))
	}
	endAt, err := utils.ParseTime(req.EndAt)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateVoucher - error parsing endAt")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidEndAt))
	}

	// update voucher
	updated, err := s.voucherRepository.UpdateVoucher(ctx, &entity.Voucher{
		ID:            id,
		Name:          req.Name,
		VoucherQuota:  req.VoucherQuota,
		Code:          req.Code,
		Discount:      req.Discount,
		StartAt:       startAt,
		EndAt:         endAt,
		VoucherTypeID: countVoucherTypeResult.ID,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrVoucherNotFound) {
			log.Error().Err(err).Msg("service::UpdateVoucher - voucher not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrVoucherNotFound))
		}

		if strings.Contains(err.Error(), constants.ErrVoucherCodeAlreadyRegistered) {
			log.Error().Err(err).Msg("service::UpdateVoucher - voucher code already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrVoucherCodeAlreadyRegistered))
		}

		log.Error().Err(err).Msg("service::UpdateVoucher - updating voucher")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response
	return &dto.CreateOrUpdateVoucherResponse{
		Name:          updated.Name,
		VoucherQuota:  updated.VoucherQuota,
		Code:          updated.Code,
		Discount:      updated.Discount,
		VoucherTypeID: countVoucherTypeResult.Type,
		CreatedAt:     utils.FormatTime(updated.CreatedAt),
		StartAt:       utils.FormatTime(updated.StartAt),
		EndAt:         utils.FormatTime(updated.EndAt),
	}, nil
}
