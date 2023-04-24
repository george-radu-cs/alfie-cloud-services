package config

import (
	"api/app/delivery"
	pb "api/app/protobuf"
	"api/app/repository"
	"api/app/services"
	"api/app/usecase"
	"api/app/utils"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"os"
)

func CreateGRPCServer(db *gorm.DB) *grpc.Server {
	repo := repository.New(db)
	mailVerifierService := services.NewMailVerifierService()
	validationsService := services.NewValidationsService()
	mediaCloudService := services.NewMediaCloudService()
	useCase := usecase.New(repo, mailVerifierService, validationsService, mediaCloudService)

	jwtService := services.NewJWTService()
	authInterceptor := delivery.NewAuthInterceptor(jwtService)

	grpcServerOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Authorize()),
	}

	if os.Getenv("IN_PRODUCTION") == "true" {
		grpcServerOptions = append(grpcServerOptions, grpc.Creds(LoadTLSCredentials()))
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	alfieServer := delivery.NewAlfieServer(useCase, jwtService)
	pb.RegisterAlfieServer(grpcServer, alfieServer)

	return grpcServer
}

func GetGRPCListener() net.Listener {
	lis, err := net.Listen("tcp", os.Getenv("GRPC_HOST")+":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		utils.ErrorLogger.Fatalf("failed to listen: %v", err)
	}
	return lis
}
