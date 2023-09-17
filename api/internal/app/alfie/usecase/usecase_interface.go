package usecase

import (
	"context"

	"api/internal/app/alfie/models"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (err error)
	VerifyUserAccount(ctx context.Context, email, code string) (err error)
	ResendUserVerificationCode(ctx context.Context, email, password string) (err error)
	Login(ctx context.Context, email, password string) (err error)
	VerifyLoginCode(ctx context.Context, email, code string) (user *models.User, err error)
	ForgotPassword(ctx context.Context, email string) (err error)
	ResetPassword(ctx context.Context, email, code, newPassword string) (err error)

	UpdateUserInfo(ctx context.Context, firstName, lastName, userEmail string) (err error)
	UpdateUserPassword(ctx context.Context, oldPassword, newPassword, userEmail string) (err error)

	CreateUploadURLForCardsDatabaseBackupForUser(ctx context.Context, userEmail string) (
		databaseUploadURL string, userDatabaseBackupFileName string, err error,
	)
	CreateDownloadURLForCardsDatabaseBackupForUser(ctx context.Context, userEmail string) (
		databaseDownloadURL string, err error,
	)
	CreateMediaFilesUploadURLsForUser(
		ctx context.Context, fileNames []string, userEmail string,
	) (filesUploadURLs []string, err error)
	CreateMediaFilesDownloadURLsForUser(ctx context.Context, requestedFileNames []string, userEmail string) (
		filesDownloadURLs []string, fileNames []string, err error,
	)
	DeleteUnusedMediaFilesForUser(ctx context.Context, activeFileNames []string, userEmail string) (err error)
}
