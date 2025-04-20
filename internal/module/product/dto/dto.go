package dto

import (
	productCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
	productSubCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	productVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/dto"
)

type CreateOrUpdateProductRequest struct {
	ProductData string `json:"product_data" validate:"required,json_string"`
}

type CreateOrUpdateProductResponse struct {
	ID             int                             `json:"id"`
	Name           string                          `json:"name"`
	Description    string                          `json:"desc"`
	Specification  string                          `json:"spec"`
	RealPrice      int                             `json:"real_price"`
	CapitalPrice   int                             `json:"capital_price"`
	DiscountPrice  int                             `json:"discount_price"`
	Stock          int                             `json:"stock"`
	Weight         float64                         `json:"weight"`
	CategoryID     int                             `json:"category_id"`
	SubCategoryID  int                             `json:"subcategory_id"`
	VariantName    string                          `json:"variant_name"`
	ProductVariant []productVariant.ProductVariant `json:"product_variant"`
	ProductGrocery []productGrocery.ProductGrocery `json:"product_grocery"`
	ProductImage   []productImage.ProductImage     `json:"product_image"`
}

type ProductData struct {
	Name          string                        `json:"name" validate:"required,min=3,max=100,xss_safe"`
	Description   string                        `json:"desc" validate:"required,max=255,xss_safe"`
	Spec          string                        `json:"spec" validate:"required,max=255,xss_safe"`
	RealPrice     int                           `json:"real_price" validate:"required,numeric,non_zero_integer,gt=0"`
	DiscountPrice int                           `json:"discount_price" validate:"numeric"`
	Stock         int                           `json:"stock" validate:"numeric"`
	CapitalPrice  int                           `json:"capital_price" validate:"numeric"`
	VariantName   string                        `json:"variant_name" validate:"required,min=3,max=100,xss_safe"`
	Variants      []productVariant.Variant      `json:"variants"`
	GroceryPrices []productGrocery.GroceryPrice `json:"grocery_prices"`
	Weight        float64                       `json:"weight" validate:"numeric"`
	CategoryID    string                        `json:"category_id" validate:"required,xss_safe"`
	SubCategoryID string                        `json:"subcategory_id" validate:"required,xss_safe"`
	ImageToKeep   []int                         `json:"image_to_keep,omitempty"`
	DelVariants   []int                         `json:"del_variants,omitempty"`
}

type UploadFileRequest struct {
	ObjectName     string `json:"object_name"`
	File           []byte `json:"-"`
	FileHeaderSize int64  `json:"-"`
	ContentType    string `json:"-"`
	Filename       string `json:"-"`
}

type GetListProductResponse struct {
	Product     []GetListProduct `json:"product"`
	TotalPages  int              `json:"total_page"`
	CurrentPage int              `json:"current_page"`
	PageSize    int              `json:"page_size"`
	TotalData   int              `json:"total_data"`
}

type GetListProduct struct {
	ID                 int                                   `json:"id"`
	Name               string                                `json:"name"`
	Description        string                                `json:"desc"`
	Specification      string                                `json:"spec"`
	RealPrice          int                                   `json:"real_price"`
	CapitalPrice       int                                   `json:"capital_price"`
	DiscountPrice      int                                   `json:"discount_price"`
	Stock              int                                   `json:"stock"`
	Weight             float64                               `json:"weight"`
	ProductCategory    productCategory.GetListCategory       `json:"product_category"`
	ProductSubCategory productSubCategory.ProductSubCategory `json:"product_subcategory"`
	VariantName        string                                `json:"variant_name"`
	ProductVariant     []productVariant.ProductVariant       `json:"product_variant"`
	ProductGrocery     []productGrocery.ProductGrocery       `json:"product_grocery"`
	ProductImage       []productImage.ProductImage           `json:"product_image"`
}
