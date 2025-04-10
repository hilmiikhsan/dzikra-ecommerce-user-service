package utils

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/dto"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeCreateOrUpdateProductResponse sanitizes the CreateOrUpdateProductResponse by removing any potentially harmful content
func SanitizeCreateOrUpdateProductResponse(resp dto.CreateOrUpdateProductResponse, policy *bluemonday.Policy) dto.CreateOrUpdateProductResponse {
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

// GetListProduct sanitizes the GetListProduct by removing any potentially harmful content
func SanitizeGetListProduct(resp dto.GetListProduct, policy *bluemonday.Policy) dto.GetListProduct {
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
