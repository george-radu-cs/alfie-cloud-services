package delivery

import (
	"api/internal/app/alfie/services"
	"api/internal/app/alfie/usecase"
	pb "api/internal/pkg/protobuf"
)

type server struct {
	pb.UnimplementedAlfieServer
	Uc         usecase.UseCase
	JWTService services.JWTService
}

func NewAlfieServer(uc usecase.UseCase, jwtService services.JWTService) pb.AlfieServer {
	return &server{Uc: uc, JWTService: jwtService}
}
