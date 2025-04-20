package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

func (s *bannerService) CreateBanner(ctx context.Context, description string, payloadFile dto.UploadFileRequest) (*dto.CreateOrUpdateBannerResponse, error) {
	// mapping file upload
	ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
	objectName := fmt.Sprintf("banner_images/%s_%s", utils.GenerateBucketFileUUID(), ext)
	byteFile := utils.NewByteFile(payloadFile.File)

	// insert new banner
	bannerResult, err := s.bannerRepository.InsertNewBanner(ctx, &entity.Banner{
		ImageURL:    objectName,
		Description: description,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::CreateBanner - Failed to insert new banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload file to minio
	uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateBanner - Failed to upload file to minio")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("Uploaded image URL: %s", uploadedPath)

	// mapping response
	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateBannerResponse{
		ID:          bannerResult.ID,
		ImageURL:    utils.FormatMediaPathURL(bannerResult.ImageURL, publicURL),
		Description: bannerResult.Description,
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateBannerResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *bannerService) GetListBanner(ctx context.Context, page, limit int, search string) (*dto.GetListBannerResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list banner
	banners, total, err := s.bannerRepository.FindListBanner(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListBanner - error getting list banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if banners is nil
	if banners == nil {
		banners = []dto.GetListBanner{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListBannerResponse{
		Banner:      banners,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *bannerService) UpdateBanner(ctx context.Context, id int, description string, payloadFile dto.UploadFileRequest) (*dto.CreateOrUpdateBannerResponse, error) {
	// find banner by id
	bannerResult, err := s.bannerRepository.FindBannerByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBannerNotFound) {
			log.Error().Err(err).Msg("service::UpdateBanner - banner not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrBannerNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateBanner - Failed to find banner by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if banner result is empty
	if bannerResult.ImageURL != "" {
		// delete old banner image
		if err := s.minioService.DeleteFile(ctx, bannerResult.ImageURL); err != nil {
			log.Error().Err(err).Msg("service::UpdateBanner - Failed to delete old banner image")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// mapping file upload
	ext := strings.ToLower(filepath.Ext(payloadFile.Filename))
	objectName := fmt.Sprintf("banner_images/%s_%s", utils.GenerateBucketFileUUID(), ext)
	byteFile := utils.NewByteFile(payloadFile.File)

	// update new banner
	bannerNewResult, err := s.bannerRepository.UpdateBanner(ctx, &entity.Banner{
		ImageURL:    objectName,
		Description: description,
		ID:          id,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::CreateBanner - Failed to update new banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// upload file to minio
	uploadedPath, err := s.minioService.UploadFile(ctx, objectName, byteFile, payloadFile.FileHeaderSize, payloadFile.ContentType)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateBanner - Failed to upload file to minio")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Msgf("Uploaded image URL: %s", uploadedPath)

	// mapping response
	publicURL := config.Envs.MinioStorage.PublicURL
	response := &dto.CreateOrUpdateBannerResponse{
		ID:          bannerNewResult.ID,
		ImageURL:    utils.FormatMediaPathURL(bannerNewResult.ImageURL, publicURL),
		Description: bannerNewResult.Description,
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateBannerResponse(*response, policy)

	return &sanitizedResponse, nil
}

func (s *bannerService) RemoveBanner(ctx context.Context, id int) error {
	// find banner by id
	bannerResult, err := s.bannerRepository.FindBannerByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrBannerNotFound) {
			log.Error().Err(err).Msg("service::RemoveBanner - banner not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrBannerNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveBanner - Failed to find banner by ID")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// delete old banner image
	if err := s.minioService.DeleteFile(ctx, bannerResult.ImageURL); err != nil {
		log.Error().Err(err).Msg("service::RemoveBanner - Failed to delete old banner image")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete banner
	if err := s.bannerRepository.SoftDeleteBannerByID(ctx, id); err != nil {
		if strings.Contains(err.Error(), constants.ErrBannerNotFound) {
			log.Error().Err(err).Msg("service::RemoveBanner - Banner not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrBannerNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveBanner - Failed to soft delete banner")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
