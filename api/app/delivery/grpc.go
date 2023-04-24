package delivery

import (
	pb "api/app/protobuf"
	"api/app/services"
	"api/app/usecase"
)

type server struct {
	pb.UnimplementedAlfieServer
	Uc         usecase.UseCase
	JWTService services.JWTService
}

func NewAlfieServer(uc usecase.UseCase, jwtService services.JWTService) pb.AlfieServer {
	return &server{Uc: uc, JWTService: jwtService}
}
