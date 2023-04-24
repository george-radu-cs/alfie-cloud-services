package usecase

import (
	"api/app/models"
)

type UseCase interface {
	Register(user *models.User) (err error)
	VerifyUserAccount(email, code string) (err error)
	ResendUserVerificationCode(email, password string) (err error)
	Login(email, password string) (err error)
	VerifyLoginCode(email, code string) (user *models.User, err error)
	ForgotPassword(email string) (err error)
	ResetPassword(email, code, newPassword string) (err error)

	UpdateUserInfo(firstName, lastName, userEmail string) (err error)
	UpdateUserPassword(oldPassword, newPassword, userEmail string) (err error)

	CreateUploadURLForCardsDatabaseBackupForUser(userEmail string) (
		databaseUploadURL string, userDatabaseBackupFileName string, err error,
	)
	CreateDownloadURLForCardsDatabaseBackupForUser(userEmail string) (databaseDownloadURL string, err error)
	CreateMediaFilesUploadURLsForUser(fileNames []string, userEmail string) (filesUploadURLs []string, err error)
	CreateMediaFilesDownloadURLsForUser(requestedFileNames []string, userEmail string) (
		filesDownloadURLs []string, fileNames []string, err error,
	)
	DeleteUnusedMediaFilesForUser(activeFileNames []string, userEmail string) (err error)
}
