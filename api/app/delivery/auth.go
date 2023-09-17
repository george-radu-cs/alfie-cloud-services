package delivery

import (
	"context"
	"errors"

	"api/app/models"
	pb "api/app/protobuf"
	"api/app/utils"
)

func (s *server) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterReply, error) {
	err := s.Uc.Register(
		ctx,
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
	ctx context.Context, request *pb.VerifyUserAccountRequest,
) (*pb.VerifyUserAccountReply, error) {
	err := s.Uc.VerifyUserAccount(ctx, request.GetEmail(), request.GetCode())
	if err != nil {
		return nil, err
	}

	return &pb.VerifyUserAccountReply{}, nil
}

func (s *server) ResendUserVerificationCode(
	ctx context.Context, request *pb.ResendUserVerificationCodeRequest,
) (*pb.ResendUserVerificationCodeReply, error) {
	err := s.Uc.ResendUserVerificationCode(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.ResendUserVerificationCodeReply{}, nil
}

func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	err := s.Uc.Login(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		// return error with code
		return nil, err
	}

	return &pb.LoginReply{}, nil
}

func (s *server) VerifyLoginCode(
	ctx context.Context, request *pb.VerifyLoginCodeRequest,
) (*pb.VerifyLoginCodeReply, error) {
	user, err := s.Uc.VerifyLoginCode(ctx, request.GetEmail(), request.GetCode())
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

func (s *server) ForgotPassword(ctx context.Context, request *pb.ForgotPasswordRequest) (
	*pb.ForgotPasswordReply, error,
) {
	err := s.Uc.ForgotPassword(ctx, request.GetEmail())
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordReply{}, nil
}

func (s *server) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordReply, error) {
	err := s.Uc.ResetPassword(ctx, request.GetEmail(), request.GetCode(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.ResetPasswordReply{}, nil
}
