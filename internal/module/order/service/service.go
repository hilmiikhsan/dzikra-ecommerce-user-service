package service

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/dto"
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *orderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest, locals *middleware.Locals, addressID, voucherID int) (*dto.CreateOrderResponse, error) {
	// convert userID to UUID
	userUUID, err := uuid.Parse(locals.UserID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - failed to parse user_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if address is exist
	addressData, err := s.addressRepository.FindDetailAddressByID(ctx, addressID, userUUID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAddressNotFound) {
			log.Error().Err(err).Msg("service::CreateOrder - address not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrAddressNotFound))
		}

		log.Error().Err(err).Msg("service::CreateOrder - failed to get detail address")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if voucher is exist
	var voucherData *entity.Voucher
	if req.VoucherID != nil {
		voucherData, err = s.voucherRepository.FindVoucherByID(ctx, voucherID)
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrVoucherNotFound) {
				log.Error().Err(err).Msg("service::CreateOrder - voucher not found")
				return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrVoucherNotFound))
			}

			log.Error().Err(err).Msg("service::CreateOrder - error finding voucher by ID")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// get list cart item
	cartItems, err := s.cartRepository.FindListCartByUserID(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - Failed to get list cart item")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if cart item is empty
	if len(cartItems) == 0 {
		log.Error().Msg("service::CreateOrder - Cart item is empty")
		return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrCartItemIsEmpty))
	}

	// sum total weight (in kg) and convert to grams
	var totalWeightKg float64
	for _, item := range cartItems {
		var weightKg float64

		if item.ProductVariantID != 0 {
			weightKg = item.ProductVariantWeight
		} else {
			weightKg = item.ProductWeight
		}

		totalWeightKg += weightKg * float64(item.Quantity)
	}

	totalWeightGrams := strconv.Itoa(int(totalWeightKg * 1000))

	// call RajaOngkir integration
	costs, err := s.rajaongkirService.GetShippingCost(ctx, totalWeightGrams, req.CostName, addressData)
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - Failed to get shipping cost")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if costs is empty
	if len(costs) == 0 || len(costs[0].Cost) == 0 {
		log.Error().Msg("service::CreateOrder - No shipping options returned")
		return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrNoShippingOptionReturned))
	}

	// check if selected service is available
	var chosenCost int
	found := false
	for _, svc := range costs[0].Cost {
		if svc.Service == req.CostService {
			chosenCost = svc.Cost[0].Value
			found = true
			break
		}
	}
	if !found {
		log.Error().Msg("service::CreateOrder - Selected service not found in shipping options")
		return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrSelectedShippingOptionNotFound))
	}

	// map product ID to struct
	prodSet := make(map[int]struct{}, len(cartItems))
	for _, it := range cartItems {
		if it.ProductVariantID == 0 {
			prodSet[it.ProductID] = struct{}{}
		}
	}

	// map product ID to grocery prices
	groceryMap := make(map[int][]productGrocery.GroceryPrice, len(prodSet))
	for pid := range prodSet {
		prices, err := s.productGroceryRepository.FindProductGroceryByProductID(ctx, pid)
		if err != nil {
			log.Error().Err(err).Msgf("service::CreateOrder - gagal load grocery for product %d", pid)
			return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		sort.Slice(prices, func(i, j int) bool {
			return prices[i].MinBuy > prices[j].MinBuy
		})

		groceryMap[pid] = prices
	}

	var (
		totalQuantity   int32
		totalProdAmount int64
	)

	computedWeights := make([]float64, len(cartItems))
	computedPrices := make([]int64, len(cartItems))

	for i, it := range cartItems {
		if it.ProductVariantID == 0 {
			for _, g := range groceryMap[it.ProductID] {
				if it.Quantity >= g.MinBuy {
					basePrice, _ := strconv.ParseInt(it.ProductRealPrice, 10, 64)
					discPrice := basePrice * int64(100-g.Discount) / 100
					cartItems[i].ProductRealPrice = fmt.Sprintf("%d", discPrice)
					break
				}
			}
		}

		var (
			weightKg  float64
			unitPrice int64
		)
		if it.ProductVariantID != 0 {
			weightKg = it.ProductVariantWeight
			vDisc, _ := strconv.ParseInt(it.ProductVariantDiscountPrice, 10, 64)
			vReal, _ := strconv.ParseInt(it.ProductVariantRealPrice, 10, 64)
			if vDisc > 0 {
				unitPrice = vDisc
			} else {
				unitPrice = vReal
			}
		} else {
			weightKg = it.ProductWeight
			pDisc, _ := strconv.ParseInt(it.ProductDiscountPrice, 10, 64)
			pReal, _ := strconv.ParseInt(it.ProductRealPrice, 10, 64)
			if pDisc > 0 {
				unitPrice = pDisc
			} else {
				unitPrice = pReal
			}
		}

		computedWeights[i] = weightKg
		computedPrices[i] = unitPrice

		totalQuantity++
		totalProdAmount += unitPrice * int64(it.Quantity)
	}

	var voucherDiscAmt int64
	if voucherData != nil {
		voucherDiscAmt = totalProdAmount * int64(voucherData.Discount) / 100
	}

	shippingCost := int64(chosenCost)
	totalShippingAmount := shippingCost
	grandTotal := totalProdAmount + shippingCost - voucherDiscAmt
	shippingAddr := fmt.Sprintf(
		"%s, %s, %s, %s%s",
		addressData.Province,
		addressData.City,
		addressData.SubDistrict,
		addressData.Address,
		func() string {
			if addressData.PostalCode != "" {
				return ", Kode Pos: " + addressData.PostalCode
			}
			return ""
		}(),
	)

	now := utils.FormatTimeJakarta()
	uuid, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - failed to generate UUID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var orderCartItems []*order.CartItem
	for _, item := range cartItems {
		var productGroceries []*order.ProductGrocery
		for _, g := range item.ProductGrocery {
			productGroceries = append(productGroceries, &order.ProductGrocery{
				Id:        int64(g.ID),
				MinBuy:    int64(g.MinBuy),
				Discount:  int64(g.Discount),
				ProductId: int64(g.ProductID),
			})
		}

		var productImages []*order.ProductImage
		for _, img := range item.ProductImage {
			productImages = append(productImages, &order.ProductImage{
				Id:        int64(img.ID),
				ImageUrl:  img.ImageURL,
				Position:  int64(img.Position),
				ProductId: int64(img.ProductID),
			})
		}

		orderCartItems = append(orderCartItems, &order.CartItem{
			Id:                          int64(item.ID),
			Quantity:                    int64(item.Quantity),
			ProductId:                   int64(item.ProductID),
			ProductVariantId:            int64(item.ProductVariantID),
			ProductName:                 item.ProductName,
			ProductRealPrice:            item.ProductRealPrice,
			ProductDiscountPrice:        item.ProductDiscountPrice,
			ProductStock:                int64(item.ProductStock),
			ProductWeight:               item.ProductWeight,
			ProductVariantWeight:        item.ProductVariantWeight,
			ProductVariantName:          item.ProductVariantName,
			ProductGroceries:            productGroceries,
			ProductVariantSubName:       item.ProductVariantSubName,
			ProductVariantRealPrice:     item.ProductVariantRealPrice,
			ProductVariantDiscountPrice: item.ProductVariantDiscountPrice,
			ProductVariantStock:         int64(item.ProductVariantStock),
			ProductImages:               productImages,
		})
	}

	newOrder := &order.CreateOrderRequest{
		Id:                  uuid.String(),
		UserId:              locals.UserID,
		Email:               locals.Email,
		AddressId:           int64(addressID),
		Status:              constants.OrderStatusUnpaid,
		ShippingName:        addressData.ReceivedName,
		ShippingAddress:     shippingAddr,
		ShippingPhone:       locals.PhoneNumber,
		ShippingType:        req.CostName,
		TotalWeight:         totalWeightKg,
		TotalQuantity:       int32(totalQuantity),
		TotalProductAmount:  totalProdAmount,
		TotalShippingCost:   shippingCost,
		TotalShippingAmount: totalShippingAmount,
		TotalAmount:         grandTotal,
		VoucherDiscount:     int64(voucherDiscAmt),
		Notes:               req.Notes,
		CostName:            req.CostName,
		CostService:         req.CostService,
		VoucherId:           int64(voucherID),
		CartItems:           orderCartItems,
		CreatedAt:           timestamppb.New(now),
	}

	res, err := s.externalOrder.CreateOrder(ctx, newOrder)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrder - Failed to create order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateOrder - Failed to rollback transaction")
			}
		}
	}()

	err = s.cartRepository.DeleteCartByUserID(ctx, tx, userUUID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrCartNotFound) {
			log.Error().Err(err).Msg("service::CreateOrder - Cart item is empty")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrCartItemIsEmpty))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to delete cart item")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateOrder - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrderResponse{
		Order: dto.OrderDetail{
			ID:                  res.Order.GetId(),
			OrderDate:           res.Order.GetOrderDate(),
			Status:              res.Order.GetStatus(),
			ShippingName:        res.Order.GetShippingName(),
			ShippingAddress:     res.Order.GetShippingAddress(),
			ShippingPhone:       res.Order.GetShippingPhone(),
			ShippingNumber:      res.Order.GetShippingNumber(),
			ShippingType:        res.Order.GetShippingType(),
			TotalWeight:         int(res.Order.GetTotalWeight()),
			TotalQuantity:       int(res.Order.GetTotalQuantity()),
			TotalShippingCost:   res.Order.GetTotalShippingCost(),
			TotalProductAmount:  res.Order.GetTotalProductAmount(),
			TotalShippingAmount: res.Order.GetTotalShippingAmount(),
			TotalAmount:         res.Order.GetTotalAmount(),
			VoucherDiscount:     int(res.Order.GetVoucherDiscount()),
			VoucherID:           res.Order.GetVoucherId(),
			CostName:            res.Order.GetCostName(),
			CostService:         res.Order.GetCostService(),
			AddressID:           int(res.Order.GetAddressId()),
			UserID:              res.Order.GetUserId(),
			Notes:               res.Order.GetNotes(),
		},
		MidtransRedirectUrl: res.GetMidtransRedirectUrl(),
	}, nil
}

func (s *orderService) GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*dto.GetListOrderResponse, error) {
	res, err := s.externalOrder.GetListOrder(ctx, page, limit, search, status, userID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrder - Failed to get list order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	orders := make([]dto.GetListOrder, 0, len(res.Orders))
	for _, order := range res.Orders {
		var voucherIDStr *string
		if order.VoucherId != 0 {
			s := fmt.Sprintf("%d", order.VoucherId)
			voucherIDStr = &s
		}

		var orderItems []dto.OrderItem
		for _, item := range order.OrderItems {
			var productImages []dto.ProductImage
			for _, img := range item.ProductImages {
				publicURL := config.Envs.MinioStorage.PublicURL
				productImages = append(productImages, dto.ProductImage{
					ID:        int(img.Id),
					ImageURL:  utils.FormatMediaPathURL(img.ImageUrl, publicURL),
					Position:  int(img.Position),
					ProductID: int(img.ProductId),
				})
			}

			var productDiscount *string
			if item.ProductDisc != 0 {
				s := fmt.Sprintf("%d", item.ProductDisc)
				productDiscount = &s
			} else {
				productDiscount = nil
			}

			orderItems = append(orderItems, dto.OrderItem{
				ProductID:             int(item.ProductId),
				ProductName:           item.ProductName,
				ProductVariantSubName: item.ProductVariantSubName,
				ProductVariant:        item.ProductVariant,
				TotalAmount:           fmt.Sprintf("%d", item.TotalAmount),
				ProductDisc:           productDiscount,
				Quantity:              int(item.Quantity),
				FixPricePerItem:       fmt.Sprintf("%d", item.FixPricePerItem),
				ProductImages:         productImages,
			})
		}

		var district *string
		if order.Address.District != "" {
			s := order.Address.District
			district = &s
		} else {
			district = nil
		}

		orders = append(orders, dto.GetListOrder{
			ID:                  order.Id,
			OrderDate:           order.OrderDate,
			Status:              order.Status,
			TotalQuantity:       int(order.TotalQuantity),
			TotalAmount:         fmt.Sprintf("%d", order.TotalAmount),
			ShippingNumber:      order.ShippingNumber,
			TotalShippingAmount: fmt.Sprintf("%d", order.TotalShippingAmount),
			CostName:            order.CostName,
			CostService:         order.CostService,
			VoucherID:           voucherIDStr,
			VoucherDiscount:     int(order.VoucherDisc),
			UserID:              userID,
			Notes:               order.Notes,
			SubTotal:            fmt.Sprintf("%d", order.SubTotal),
			Address: dto.Address{
				ID:                  int(order.Address.Id),
				Province:            order.Address.Province,
				City:                order.Address.City,
				District:            district,
				SubDistrict:         order.Address.Subdistrict,
				PostalCode:          order.Address.PostalCode,
				Address:             order.Address.Address,
				ReceivedName:        order.Address.ReceivedName,
				UserID:              userID,
				CityVendorID:        order.Address.CityVendorId,
				ProvinceVendorID:    order.Address.ProvinceVendorId,
				SubDistrictVendorID: order.Address.SubdistrictVendorId,
			},
			OrderItems: orderItems,
			Payment: dto.Payment{
				RedirectURL: order.Payment.RedirectUrl,
			},
		})
	}

	return &dto.GetListOrderResponse{
		Orders:      orders,
		TotalPages:  int(res.TotalPages),
		CurrentPage: int(res.CurrentPage),
		PageSize:    int(res.PageSize),
		TotalData:   int(res.TotalData),
	}, nil
}

func (s *orderService) GetWaybillDetails(ctx context.Context, orderID string) (*dto.GetWaybillResponse, error) {
	res, err := s.externalOrder.GetOrderById(ctx, &order.GetOrderByIdRequest{
		Id: orderID,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrOrderNotFound) {
			log.Error().Err(err).Msg("service::GetWaybillDetails - Order not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		log.Error().Err(err).Msg("service::GetWaybillDetails - Failed to get order details")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if res == nil {
		log.Error().Msg("service::GetWaybillDetails - Order not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
	}
	if res.Order.ShippingName == "" || res.Order.ShippingType == "" {
		log.Error().Msg("service::GetWaybillDetails - Shipping details not available for this order")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("shipping details not available for this order"))
	}

	waybill, err := s.rajaongkirService.GetWaybill(ctx, res.Order.ShippingNumber, strings.ToLower(res.Order.ShippingType))
	if err != nil {
		log.Error().Err(err).Msg("service::GetWaybillDetails - Failed to get waybill details")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	result := &dto.GetWaybillResponse{
		Delivered:    waybill.DeliveryStatus.PODReceiver != "",
		Destination:  waybill.Summary.Destination,
		Resi:         waybill.Summary.Resi,
		ServiceCode:  waybill.Summary.ServiceCode,
		WaybillDate:  waybill.Summary.WaybillDate,
		ShipperName:  waybill.Summary.ShipperName,
		ReceiverName: waybill.Summary.ReceiverName,
		Origin:       waybill.Summary.Origin,
		Status:       waybill.Summary.Status,
		CourierName:  waybill.Summary.CourierName,
		Manifest:     []dto.WaybillManifest{},
		DeliveryStatus: dto.DeliveryStatus{
			Status:      waybill.DeliveryStatus.Status,
			PODReceiver: waybill.DeliveryStatus.PODReceiver,
			PODDate:     waybill.DeliveryStatus.PODDate,
			PODTime:     waybill.DeliveryStatus.PODTime,
		},
	}

	for _, m := range waybill.Manifest {
		result.Manifest = append(result.Manifest, dto.WaybillManifest{
			Description: m.Description,
			Date:        m.Date,
			Time:        m.Time,
			City:        m.CityName,
		})
	}

	return result, nil
}

func (s *orderService) GetListOrderTransaction(ctx context.Context, page, limit int, search, status string) (*dto.GetListOrderResponse, error) {
	res, err := s.externalOrder.GetListOrderTransaction(ctx, page, limit, search, status)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListOrderTransaction - Failed to get list order")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	orders := make([]dto.GetListOrder, 0, len(res.Orders))
	for _, order := range res.Orders {
		var voucherIDStr *string
		if order.VoucherId != 0 {
			s := fmt.Sprintf("%d", order.VoucherId)
			voucherIDStr = &s
		}

		var orderItems []dto.OrderItem
		for _, item := range order.OrderItems {
			var productImages []dto.ProductImage
			for _, img := range item.ProductImages {
				publicURL := config.Envs.MinioStorage.PublicURL
				productImages = append(productImages, dto.ProductImage{
					ID:        int(img.Id),
					ImageURL:  utils.FormatMediaPathURL(img.ImageUrl, publicURL),
					Position:  int(img.Position),
					ProductID: int(img.ProductId),
				})
			}

			var productDiscount *string
			if item.ProductDisc != 0 {
				s := fmt.Sprintf("%d", item.ProductDisc)
				productDiscount = &s
			} else {
				productDiscount = nil
			}

			orderItems = append(orderItems, dto.OrderItem{
				ProductID:             int(item.ProductId),
				ProductName:           item.ProductName,
				ProductVariantSubName: item.ProductVariantSubName,
				ProductVariant:        item.ProductVariant,
				TotalAmount:           fmt.Sprintf("%d", item.TotalAmount),
				ProductDisc:           productDiscount,
				Quantity:              int(item.Quantity),
				FixPricePerItem:       fmt.Sprintf("%d", item.FixPricePerItem),
				ProductImages:         productImages,
			})
		}

		var district *string
		if order.Address.District != "" {
			s := order.Address.District
			district = &s
		} else {
			district = nil
		}

		orders = append(orders, dto.GetListOrder{
			ID:                  order.Id,
			OrderDate:           order.OrderDate,
			Status:              order.Status,
			TotalQuantity:       int(order.TotalQuantity),
			TotalAmount:         fmt.Sprintf("%d", order.TotalAmount),
			ShippingNumber:      order.ShippingNumber,
			TotalShippingAmount: fmt.Sprintf("%d", order.TotalShippingAmount),
			CostName:            order.CostName,
			CostService:         order.CostService,
			VoucherID:           voucherIDStr,
			VoucherDiscount:     int(order.VoucherDisc),
			UserID:              order.Address.UserId,
			Notes:               order.Notes,
			SubTotal:            fmt.Sprintf("%d", order.SubTotal),
			Address: dto.Address{
				ID:                  int(order.Address.Id),
				Province:            order.Address.Province,
				City:                order.Address.City,
				District:            district,
				SubDistrict:         order.Address.Subdistrict,
				PostalCode:          order.Address.PostalCode,
				Address:             order.Address.Address,
				ReceivedName:        order.Address.ReceivedName,
				UserID:              order.Address.UserId,
				CityVendorID:        order.Address.CityVendorId,
				ProvinceVendorID:    order.Address.ProvinceVendorId,
				SubDistrictVendorID: order.Address.SubdistrictVendorId,
			},
			OrderItems: orderItems,
			Payment: dto.Payment{
				RedirectURL: order.Payment.RedirectUrl,
			},
		})
	}

	return &dto.GetListOrderResponse{
		Orders:      orders,
		TotalPages:  int(res.TotalPages),
		CurrentPage: int(res.CurrentPage),
		PageSize:    int(res.PageSize),
		TotalData:   int(res.TotalData),
	}, nil
}

func (s *orderService) UpdateOrderShippingNumber(ctx context.Context, req *dto.UpdateOrderShippingNumberRequest, orderID string) (*dto.UpdateOrderShippingNumberResponse, error) {
	res, err := s.externalOrder.UpdateOrderShippingNumber(ctx, &order.UpdateOrderShippingNumberRequest{
		Id:             orderID,
		ShippingNumber: req.ShippingNumber,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrOrderNotFound) {
			log.Error().Err(err).Msg("service::UpdateOrderShippingNumber - Order not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		if strings.Contains(err.Error(), constants.ErrShippingNumberAlreadyExists) {
			log.Error().Err(err).Msg("service::UpdateOrderShippingNumber - Shipping number already exists")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrShippingNumberAlreadyExists))
		}

		log.Error().Err(err).Msg("service::UpdateOrderShippingNumber - Failed to update order shipping number")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.UpdateOrderShippingNumberResponse{
		ID:                  res.Id,
		OrderDate:           res.OrderDate,
		Status:              res.Status,
		ShippingName:        res.ShippingName,
		ShippingAddress:     res.ShippingAddress,
		ShippingPhone:       res.ShippingPhone,
		ShippingNumber:      res.ShippingNumber,
		ShippingType:        res.ShippingType,
		TotalWeight:         int(res.TotalWeight),
		TotalQuantity:       int(res.TotalQuantity),
		TotalShippingCost:   res.TotalShippingCost,
		TotalProductAmount:  res.TotalProductAmount,
		TotalShippingAmount: res.TotalShippingAmount,
		TotalAmount:         res.TotalAmount,
		VoucherDiscount:     int(res.VoucherDiscount),
		VoucherID:           &res.VoucherId,
		CostName:            res.CostName,
		CostService:         res.CostService,
		AddressID:           int(res.AddressId),
		UserID:              res.UserId,
		Notes:               res.Notes,
	}, nil
}

func (s *orderService) UpdateOrderStatusTransaction(ctx context.Context, req *dto.UpdateOrderStatusTransactionRequest, orderID, fullName, email string) error {
	orderResult, err := s.externalOrder.GetOrderById(ctx, &order.GetOrderByIdRequest{
		Id: orderID,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrOrderNotFound) {
			log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Order not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to get order details")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	userFcmToken, err := s.userFcmTokenRepository.FindUserFCMTokenByUserID(ctx, orderResult.Order.UserId)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to get user FCM token")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateOrderStatusTransaction - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateOrderStatusTransaction - Failed to rollback transaction")
			}
		}
	}()

	if req.Status == constants.OrderStatusWaitingForPickup {
		orderItems, err := s.externalOrder.GetOrderItemsByOrderID(ctx, &order.GetOrderItemsByOrderIDRequest{
			OrderId: orderID,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to get order items")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		for _, item := range orderItems.OrderItems {
			var stock int
			if item.ProductVariantId != 0 {
				stock, err = s.productVariantRepository.FindProductVariantStockByID(ctx, int(item.ProductVariantId))
			} else {
				stock, err = s.productRepository.FindProductStockByID(ctx, int(item.ProductId))
			}
			if err != nil {
				log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to fetch stock")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
			if stock < int(item.Quantity) {
				log.Error().Msg("service::UpdateOrderStatusTransaction - Insufficient stock for some items in the order")
				return err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage("Stok tidak cukup untuk beberapa item dalam pesanan"))
			}
		}

		for _, item := range orderItems.OrderItems {
			if err = s.productRepository.ReduceStock(ctx, tx, int(item.ProductId), int(item.ProductVariantId), int(item.Quantity)); err != nil {
				log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to reduce stock")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}

		notifFcmRequest := &notification.SendFcmNotificationRequest{
			FcmToken:        userFcmToken,
			Title:           "Pesanan menunggu jasa kirim",
			Body:            fmt.Sprintf("Pesanan dengan id %s siap dikirim.", orderResult.Order.Id),
			UserId:          orderResult.Order.UserId,
			IsStatusChanged: true,
		}

		if _, err := s.externalNotification.SendFcmNotification(ctx, notifFcmRequest); err != nil {
			log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to send FCM notification")
			return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	_, err = s.externalOrder.UpdateOrderStatusTransaction(ctx, &order.UpdateOrderStatusTransactionRequest{
		Status:  req.Status,
		OrderId: orderID,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateOrderStatusTransaction - Failed to update order status")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateOrderStatusTransaction - Failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
