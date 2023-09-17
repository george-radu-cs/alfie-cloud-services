package delivery

import (
	pb "api/app/protobuf"
	"context"
)

func (s *server) CreateUploadURLForCardsDatabaseBackup(
	ctx context.Context, request *pb.CreateUploadURLForCardsDatabaseBackupRequest,
) (*pb.CreateUploadURLForCardsDatabaseBackupReply, error) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	databaseUploadURL, databaseFileName, err := s.Uc.CreateUploadURLForCardsDatabaseBackupForUser(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUploadURLForCardsDatabaseBackupReply{
		DatabaseUploadURL: databaseUploadURL,
		DatabaseFileName:  databaseFileName,
	}, nil
}

func (s *server) CreateDownloadURLForCardsDatabaseBackup(
	ctx context.Context, request *pb.CreateDownloadURLForCardDatabaseBackupRequest,
) (*pb.CreateDownloadURLForCardDatabaseBackupReply, error) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	databaseDownloadURL, err := s.Uc.CreateDownloadURLForCardsDatabaseBackupForUser(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDownloadURLForCardDatabaseBackupReply{
		DatabaseDownloadURL: databaseDownloadURL,
	}, nil
}

func (s *server) CreateMediaFilesUploadURLs(
	ctx context.Context, request *pb.CreateMediaFilesUploadURLsRequest,
) (*pb.CreateMediaFilesUploadURLsReply, error) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	mediaFilesUploadURLs, err := s.Uc.CreateMediaFilesUploadURLsForUser(ctx, request.GetFileNames(), userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.CreateMediaFilesUploadURLsReply{
		MediaFilesUploadURLs: mediaFilesUploadURLs,
	}, nil
}

func (s *server) CreateMediaFilesDownloadURLs(
	ctx context.Context, request *pb.CreateMediaFilesDownloadURLsRequest,
) (*pb.CreateMediaFilesDownloadURLsReply, error) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	mediaFilesDownloadURLs, mediaFileNames, err := s.Uc.CreateMediaFilesDownloadURLsForUser(
		ctx, request.GetFileNames(), userEmail,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateMediaFilesDownloadURLsReply{
		MediaFilesDownloadURLs: mediaFilesDownloadURLs,
		MediaFilesNames:        mediaFileNames,
	}, nil
}

func (s *server) DeleteUnusedMediaFiles(
	ctx context.Context, request *pb.DeleteUnusedMediaFilesRequest,
) (*pb.DeleteUnusedMediaFilesReply, error) {
	userEmail := getUserEmailFromValidatedContext(ctx)
	err := s.Uc.DeleteUnusedMediaFilesForUser(ctx, request.GetActiveFileNames(), userEmail)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUnusedMediaFilesReply{}, nil
}
