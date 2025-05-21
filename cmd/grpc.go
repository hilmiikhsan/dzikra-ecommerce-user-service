package cmd

import (
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/address"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/cart"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_image"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/product_variant"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/tokenvalidation"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	addressGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/handler/grpc"
	cartGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/handler/grpc"
	productGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/handler/grpc"
	productImageGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/handler/grpc"
	productVariantGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/handler/grpc"
	userGrpcHandler "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user/handler/grpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func RunServeGRPC() {
	envs := config.Envs

	logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	infrastructure.InitializeLogger(
		envs.App.Environtment,
		envs.App.LogFile,
		logLevel,
	)

	lis, err := net.Listen("tcp", ":"+envs.App.GrpcPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on gRPC port")
	}

	grpcServer := grpc.NewServer()

	opts := []adapter.Option{
		adapter.WithGRPCServer(grpcServer),
		adapter.WithDzikraPostgres(),
		adapter.WithDzikraRedis(),
	}

	// initialize Minio and capture any error
	minioOpt, err := adapter.WithDzikraMinio()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Minio adapter")
	}
	opts = append(opts, minioOpt)

	// Sync adapters
	if err := adapter.Adapters.Sync(opts...); err != nil {
		log.Fatal().Err(err).Msg("Failed to sync adapters")
	}

	tokenvalidation.RegisterTokenValidationServer(grpcServer, userGrpcHandler.NewUserGrpcAPI())
	cart.RegisterCartServiceServer(grpcServer, cartGrpcHandler.NewCartGrpcAPI())
	product_image.RegisterProductImageServiceServer(grpcServer, productImageGrpcHandler.NewProductImageGrpcAPI())
	address.RegisterAddressServiceServer(grpcServer, addressGrpcHandler.NewAddressGrpcAPI())
	product_variant.RegisterProductVariantServiceServer(grpcServer, productVariantGrpcHandler.NewProductVariantGrpcAPI())
	product.RegisterProductServiceServer(grpcServer, productGrpcHandler.NewProductGrpcAPI())

	go func() {
		log.Info().Msgf("gRPC server is running on port %s", envs.App.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		signals = []os.Signal{os.Interrupt}
	}
	signal.Notify(quit, signals...)
	<-quit

	log.Info().Msg("gRPC server is shutting down ...")
	grpcServer.GracefulStop()

	if err := adapter.Adapters.Unsync(); err != nil {
		log.Error().Err(err).Msg("Error while closing adapters")
	}

	log.Info().Msg("gRPC server gracefully stopped")
}
