package services

import "context"

const (
	BackupFolderName = "backup"
	MarketFolderName = "market"
	MediaFolderName  = "media"
)

type MediaCloudService interface {
	CreatePresignedURLForFileUpload(ctx context.Context, fileNameKey string) (presignedUploadURL string, err error)
	CreatePresignedURLForFileDownload(ctx context.Context, fileNameKey string) (presignedDownloadURL string, err error)
	DeleteFile(ctx context.Context, fileNameKey string) (err error)
	CreatePresignedURLsForMultipleFilesUpload(ctx context.Context, fileNameKeys []string) (
		presignedUploadURLs []string, err error,
	)
	CreatePresignedURLsForMultipleFilesDownload(
		ctx context.Context, fileNameKeys []string,
	) (presignedDownloadURLs []string, err error)
	DeleteMultipleFiles(ctx context.Context, fileNameKeys []string) (err error)
	CheckIfFileExists(ctx context.Context, fileNameKey string) (fileExists bool, err error)
	CreateFolder(ctx context.Context, folderName string) (err error)
	GetListOfFilesFromFolder(ctx context.Context, folderName string, maxNumberOfFiles int32) (
		fileKeys []string, err error,
	)
}
