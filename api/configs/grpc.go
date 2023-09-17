package configs

import (
	"net"
	"os"

	delivery2 "api/internal/app/alfie/delivery"
	"api/internal/app/alfie/repository"
	services2 "api/internal/app/alfie/services"
	"api/internal/app/alfie/usecase"
	"api/internal/app/alfie/utils"
	pb "api/internal/pkg/protobuf"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func CreateGRPCServer(db *gorm.DB) *grpc.Server {
	repo := repository.New(db)
	mailVerifierService := services2.NewMailVerifierService()
	validationsService := services2.NewValidationsService()
	mediaCloudService := services2.NewMediaCloudService()
	useCase := usecase.New(repo, mailVerifierService, validationsService, mediaCloudService)

	jwtService := services2.NewJWTService()
	authInterceptor := delivery2.NewAuthInterceptor(jwtService)

	grpcServerOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Authorize()),
	}

	if os.Getenv("IN_PRODUCTION") == "true" {
		grpcServerOptions = append(grpcServerOptions, grpc.Creds(LoadTLSCredentials()))
	}

	grpcServer := grpc.NewServer(grpcServerOptions...)
	alfieServer := delivery2.NewAlfieServer(useCase, jwtService)
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
