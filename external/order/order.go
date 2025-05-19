package order

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("ORDER_GRPC_HOST", config.Envs.Order.OrderGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::CreateOrder - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := order.NewOrderServiceClient(conn)

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::CreateOrder - Failed to create order")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::CreateOrder - Response error from order")
		return nil, fmt.Errorf("get response error from order: %s", resp.Message)
	}

	return resp, nil
}

func (*External) GetListOrder(ctx context.Context, page, limit int, search, status, userID string) (*order.GetListOrderResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("ORDER_GRPC_HOST", config.Envs.Order.OrderGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetListOrder - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := order.NewOrderServiceClient(conn)

	req := &order.GetListOrderRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: search,
		Status: status,
		UserId: userID,
	}

	resp, err := client.GetListOrder(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetListOrder - Failed to get list order")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetListOrder - Response error from order")
		return nil, fmt.Errorf("get response error from order: %s", resp.Message)
	}

	return resp, nil
}

func (*External) GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("ORDER_GRPC_HOST", config.Envs.Order.OrderGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetOrderById - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := order.NewOrderServiceClient(conn)

	resp, err := client.GetOrderById(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetOrderById - Failed to get detail order")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		if resp.Message == constants.ErrOrderNotFound {
			log.Err(err).Msg("external::GetOrderById - Order not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrOrderNotFound))
		}

		log.Err(err).Msg("external::GetOrderById - Response error from order")
		return nil, fmt.Errorf("get response error from order: %s", resp.Message)
	}

	return resp, nil
}
