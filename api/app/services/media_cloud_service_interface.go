package services

const (
	BackupFolderName = "backup"
	MarketFolderName = "market"
	MediaFolderName  = "media"
)

type MediaCloudService interface {
	CreatePresignedURLForFileUpload(fileNameKey string) (presignedUploadURL string, err error)
	CreatePresignedURLForFileDownload(fileNameKey string) (presignedDownloadURL string, err error)
	DeleteFile(fileNameKey string) (err error)
	CreatePresignedURLsForMultipleFilesUpload(fileNameKeys []string) (presignedUploadURLs []string, err error)
	CreatePresignedURLsForMultipleFilesDownload(fileNameKeys []string) (presignedDownloadURLs []string, err error)
	DeleteMultipleFiles(fileNameKeys []string) (err error)
	CheckIfFileExists(fileNameKey string) (fileExists bool, err error)
	CreateFolder(folderName string) (err error)
	GetListOfFilesFromFolder(folderName string, maxNumberOfFiles int32) (fileKeys []string, err error)
}
