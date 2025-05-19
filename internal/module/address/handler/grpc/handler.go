package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/address"
)

func (h *addressImageGrpcAPI) GetAddressesByIds(ctx context.Context, req *address.GetAddressesByIdsRequest) (*address.GetAddressesResponse, error) {
	addresses, err := h.AddressService.GetAddressesByIds(ctx, req.Ids)
	if err != nil {
		return &address.GetAddressesResponse{
			Message:   "failed to get addresses by IDs",
			Addresses: nil,
		}, nil
	}

	resp := &address.GetAddressesResponse{
		Message: "success",
	}

	for _, addr := range addresses {
		resp.Addresses = append(resp.Addresses, &address.Address{
			Id:                  int32(addr.ID),
			UserId:              addr.UserID,
			Address:             addr.Address,
			City:                addr.City,
			Province:            addr.Province,
			PostalCode:          addr.PostalCode,
			ProvinceVendorId:    addr.ProvinceVendorID,
			CityVendorId:        addr.CityVendorID,
			Subdistrict:         addr.SubDistrict,
			SubdistrictVendorId: addr.SubDistrictVendorID,
			ReceivedName:        addr.ReceivedName,
		})
	}

	return resp, nil
}
