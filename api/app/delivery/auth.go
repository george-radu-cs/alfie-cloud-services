package delivery

import (
	"api/app/models"
	pb "api/app/protobuf"
	"api/app/utils"
	"context"
	"errors"
)

func (s *server) Register(_ context.Context, request *pb.RegisterRequest) (*pb.RegisterReply, error) {
	err := s.Uc.Register(
		&models.User{
			Email:     request.GetEmail(),
			Password:  request.GetPassword(),
			FirstName: request.GetFirstName(),
			LastName:  request.GetLastName(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{}, nil
}

func (s *server) VerifyUserAccount(
	_ context.Context, request *pb.VerifyUserAccountRequest,
) (*pb.VerifyUserAccountReply, error) {
	err := s.Uc.VerifyUserAccount(request.GetEmail(), request.GetCode())
	if err != nil {
		return nil, err
	}

	return &pb.VerifyUserAccountReply{}, nil
}

func (s *server) ResendUserVerificationCode(
	_ context.Context, request *pb.ResendUserVerificationCodeRequest,
) (*pb.ResendUserVerificationCodeReply, error) {
	err := s.Uc.ResendUserVerificationCode(request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.ResendUserVerificationCodeReply{}, nil
}

func (s *server) Login(_ context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	err := s.Uc.Login(request.GetEmail(), request.GetPassword())
	if err != nil {
		// return error with code
		return nil, err
	}

	return &pb.LoginReply{}, nil
}

func (s *server) VerifyLoginCode(
	_ context.Context, request *pb.VerifyLoginCodeRequest,
) (*pb.VerifyLoginCodeReply, error) {
	user, err := s.Uc.VerifyLoginCode(request.GetEmail(), request.GetCode())
	if err != nil {
		return nil, err
	}

	token, err := s.JWTService.GenerateToken(request.GetEmail())
	if err != nil {
		utils.ErrorLogger.Printf(
			"Error generating token for user with email %s, error %s", request.GetEmail(), err.Error(),
		)
		return nil, errors.New("error generating token")
	}

	return &pb.VerifyLoginCodeReply{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     token,
	}, nil
}

func (s *server) ForgotPassword(_ context.Context, request *pb.ForgotPasswordRequest) (
	*pb.ForgotPasswordReply, error,
) {
	err := s.Uc.ForgotPassword(request.GetEmail())
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordReply{}, nil
}

func (s *server) ResetPassword(_ context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordReply, error) {
	err := s.Uc.ResetPassword(request.GetEmail(), request.GetCode(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordReply{}, nil
}
