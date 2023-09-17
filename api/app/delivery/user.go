package delivery

import (
	pb "api/app/protobuf"
	"context"
)

func (s *server) UpdateUserInfo(ctx context.Context, request *pb.UpdateUserInfoRequest) (
	*pb.UpdateUserInfoReply, error,
) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	err := s.Uc.UpdateUserInfo(ctx, request.GetFirstName(), request.GetLastName(), userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserInfoReply{}, nil
}

func (s *server) UpdatePassword(ctx context.Context, request *pb.UpdatePasswordRequest) (
	*pb.UpdatePasswordReply, error,
) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	err := s.Uc.UpdateUserPassword(ctx, request.GetOldPassword(), request.GetNewPassword(), userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.UpdatePasswordReply{}, nil
}
