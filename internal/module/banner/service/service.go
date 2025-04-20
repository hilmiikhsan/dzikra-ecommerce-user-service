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

func (s *bannerService) CreateBanner(ctx context.Context, description string, payloadFile dto.UploadFileRequest) (*dto.CreateBannerResponse, error) {
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
	response := &dto.CreateBannerResponse{
		ID:          bannerResult.ID,
		ImageURL:    utils.FormatMediaPathURL(bannerResult.ImageURL, publicURL),
		Description: bannerResult.Description,
	}

	// Sanitize response
	policy := bluemonday.UGCPolicy()
	sanitizedResponse := utils.SanitizeCreateOrUpdateBannerResponse(*response, policy)

	return &sanitizedResponse, nil
}
