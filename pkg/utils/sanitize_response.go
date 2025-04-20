package utils

import (
	banner "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/dto"
	product "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeCreateOrUpdateProductResponse sanitizes the CreateOrUpdateProductResponse by removing any potentially harmful content
func SanitizeCreateOrUpdateProductResponse(resp product.CreateOrUpdateProductResponse, policy *bluemonday.Policy) product.CreateOrUpdateProductResponse {
	resp.Name = policy.Sanitize(resp.Name)
	resp.Description = policy.Sanitize(resp.Description)
	resp.Specification = policy.Sanitize(resp.Specification)
	resp.VariantName = policy.Sanitize(resp.VariantName)

	for i, variant := range resp.ProductVariant {
		variant.VariantSubName = policy.Sanitize(variant.VariantSubName)
		resp.ProductVariant[i] = variant
	}

	for i, grocery := range resp.ProductGrocery {
		resp.ProductGrocery[i] = grocery
	}

	for i, image := range resp.ProductImage {
		image.ImageURL = policy.Sanitize(image.ImageURL)
		resp.ProductImage[i] = image
	}

	return resp
}

// SanitizeGetListProduct sanitizes the GetListProduct by removing any potentially harmful content
func SanitizeGetListProduct(resp product.GetListProduct, policy *bluemonday.Policy) product.GetListProduct {
	resp.Name = policy.Sanitize(resp.Name)
	resp.Description = policy.Sanitize(resp.Description)
	resp.Specification = policy.Sanitize(resp.Specification)
	resp.VariantName = policy.Sanitize(resp.VariantName)

	for i, variant := range resp.ProductVariant {
		variant.VariantSubName = policy.Sanitize(variant.VariantSubName)
		resp.ProductVariant[i] = variant
	}

	for i, grocery := range resp.ProductGrocery {
		resp.ProductGrocery[i] = grocery
	}

	for i, image := range resp.ProductImage {
		image.ImageURL = policy.Sanitize(image.ImageURL)
		resp.ProductImage[i] = image
	}

	return resp
}

// SanitizeCreateOrUpdateBannerResponse sanitizes the CreateOrUpdateBannerResponse by removing any potentially harmful content
func SanitizeCreateOrUpdateBannerResponse(resp banner.CreateOrUpdateBannerResponse, policy *bluemonday.Policy) banner.CreateOrUpdateBannerResponse {
	resp.ImageURL = policy.Sanitize(resp.ImageURL)
	resp.Description = policy.Sanitize(resp.Description)

	return resp
}
