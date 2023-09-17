package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"api/internal/app/alfie/models"
	"api/internal/app/alfie/services"
	"api/internal/app/alfie/utils"
)

func (uc *useCase) CreateUploadURLForCardsDatabaseBackupForUser(ctx context.Context, userEmail string) (
	databaseUploadURL string, userDatabaseBackupFileName string, err error,
) {
	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return "", "", errors.New("invalid request")
	}

	userDatabaseBackupFileName = uc.getDatabaseBackupFileNameForUser(user)
	userDatabaseBackupFileNameKey := uc.getDatabaseBackupFileNameKeyForUser(user)
	databaseUploadURL, err = uc.MediaCloudService.CreatePresignedURLForFileUpload(ctx, userDatabaseBackupFileNameKey)
	if err != nil {
		utils.ErrorLogger.Printf(
			"Error creating presigned URL to upload database backup file for user %s: %v",
			user.Email, err,
		)
		return "", "", errors.New("error processing request")
	}

	return databaseUploadURL, userDatabaseBackupFileName, nil
}

func (uc *useCase) CreateDownloadURLForCardsDatabaseBackupForUser(ctx context.Context, userEmail string) (
	databaseDownloadURL string, err error,
) {
	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return "", errors.New("invalid request")
	}

	userDatabaseBackupFileNameKey := uc.getDatabaseBackupFileNameKeyForUser(user)
	backupFileExists, err := uc.MediaCloudService.CheckIfFileExists(ctx, userDatabaseBackupFileNameKey)
	if err != nil || !backupFileExists {
		utils.ErrorLogger.Printf(
			"Error creating presigned URL to download database backup file for user %s: %v", user.Email, err,
		)
		return "", errors.New("backup file doesn't exist in cloud")
	}

	databaseDownloadURL, err = uc.MediaCloudService.CreatePresignedURLForFileDownload(
		ctx,
		userDatabaseBackupFileNameKey,
	)
	if err != nil {
		utils.ErrorLogger.Printf(
			"Error creating presigned URL to download database backup file for user %s: %v",
			user.Email, err,
		)
		return "", errors.New("error processing request")
	}

	return databaseDownloadURL, nil
}

func (uc *useCase) CreateMediaFilesUploadURLsForUser(ctx context.Context, fileNames []string, userEmail string) (
	filesUploadURLs []string, err error,
) {
	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return nil, errors.New("invalid request")
	}

	fileNameKeys := uc.createMediaFileKeysForUser(user, fileNames)
	fileUploadURLs, err := uc.MediaCloudService.CreatePresignedURLsForMultipleFilesUpload(ctx, fileNameKeys)
	if err != nil {
		utils.ErrorLogger.Printf("Error creating presigned URLs to upload media files for user %s: %v", user.Email, err)
		return nil, errors.New("error processing request")
	}

	return fileUploadURLs, nil
}

func (uc *useCase) CreateMediaFilesDownloadURLsForUser(
	ctx context.Context,
	requestedFileNames []string, userEmail string,
) (filesDownloadURLs []string, fileNames []string, err error) {
	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return nil, nil, errors.New("invalid request")
	}

	fileKeysThatWereBackedUp, fileNamesThatWereBackedUp, err := uc.filterMediaFilesForUserThatWereBackedUpInCloud(
		ctx, requestedFileNames, user,
	)
	fileDownloadURLs, err := uc.MediaCloudService.CreatePresignedURLsForMultipleFilesDownload(
		ctx, fileKeysThatWereBackedUp,
	)
	if err != nil {
		utils.ErrorLogger.Printf(
			"Error creating presigned URLs to download media files for user %s: %v", user.Email, err,
		)
		return nil, nil, errors.New("error processing request")
	}

	return fileDownloadURLs, fileNamesThatWereBackedUp, nil
}

func (uc *useCase) DeleteUnusedMediaFilesForUser(
	ctx context.Context, activeFileNames []string, userEmail string,
) (err error) {
	user, err := uc.Repository.GetUserByEmail(userEmail)
	if uc.didNotReceivedUser(err, user) {
		utils.ErrorLogger.Printf("User not found for email %s", userEmail)
		return errors.New("invalid request")
	}

	userMediaFolderKey := uc.getMediaFolderKeyForUser(user)
	userMediaFilesKeys, err := uc.MediaCloudService.GetListOfFilesFromFolder(
		ctx, userMediaFolderKey, user.S3MaxNumberOfMediaFiles,
	)
	if err != nil {
		utils.ErrorLogger.Printf("Error getting list of media files to delete for user %s: %v", user.Email, err)
		return errors.New("error processing request")
	}

	activeFileNamesKeys := uc.getMediaFileKeysForUserFromFileNames(user, activeFileNames)
	fileKeysToDelete := uc.getUnusedFileKeys(userMediaFilesKeys, activeFileNamesKeys)
	if len(fileKeysToDelete) == 0 {
		return nil
	}

	err = uc.MediaCloudService.DeleteMultipleFiles(ctx, fileKeysToDelete)
	if err != nil {
		utils.ErrorLogger.Printf("Error deleting unused media files for user %s: %v", user.Email, err)
		return errors.New("error processing request")
	}

	return nil
}

func (uc *useCase) getDatabaseBackupFileNameForUser(user *models.User) (userDatabaseName string) {
	return user.S3ID + ".db"
}

func (uc *useCase) getDatabaseBackupFileNameKeyForUser(user *models.User) (userDatabaseKey string) {
	return services.BackupFolderName + "/" + user.S3ID + ".db"
}

func (uc *useCase) getMediaFolderKeyForUser(user *models.User) (userMediaFolderKey string) {
	return services.MediaFolderName + "/" + user.S3ID + "/"
}

func (uc *useCase) getMediaFileKeyForUserFromFileName(user *models.User, fileName string) (mediaFileKey string) {
	return uc.getMediaFolderKeyForUser(user) + fileName
}

func (uc *useCase) getMediaFileKeysForUserFromFileNames(
	user *models.User, fileNames []string,
) (mediaFileKeys []string) {
	fileKeys := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		fileKeys[i] = uc.getMediaFileKeyForUserFromFileName(user, fileName)
	}
	return fileKeys
}

func (uc *useCase) createMediaFolderForUser(ctx context.Context, user *models.User) (err error) {
	userMediaFolderKey := uc.getMediaFolderKeyForUser(user)
	return uc.MediaCloudService.CreateFolder(ctx, userMediaFolderKey)
}

func (uc *useCase) createMediaFileKeysForUser(user *models.User, fileNames []string) (mediaFileKeys []string) {
	userMediaFolderKey := uc.getMediaFolderKeyForUser(user)
	fileKeys := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		fileKeys[i] = userMediaFolderKey + fileName
	}
	return fileKeys
}

func (uc *useCase) getUnusedFileKeys(userMediaFilesKeys, activeFileNames []string) (unusedFileKeys []string) {
	// the function will sort the activeFileNames strings for binary search for faster lookup
	// assume len(userMediaFilesKeys) = m & len(activeFileNames) = n
	// sorting will take O(n * log(n)) time only once
	// binary search will take O(m * log(n)) time
	// alternative for O(m * n) time is to use nested for loop for linear search

	unusedFileKeys = make([]string, 0)
	sort.Strings(activeFileNames)
	for _, userMediaFileKey := range userMediaFilesKeys {
		i := sort.SearchStrings(activeFileNames, userMediaFileKey)
		if i == len(activeFileNames) || activeFileNames[i] != userMediaFileKey {
			unusedFileKeys = append(unusedFileKeys, userMediaFileKey)
		}
	}
	return unusedFileKeys
}

// filterMediaFilesForUserThatWereBackedUpInCloud returns the fileKeys alongside the fileNames from the
// requestedFileNames that were backed up in the cloud for the user, since the user could have deleted
// some files from the cloud or the user could have requested to download files that were not backed up
func (uc *useCase) filterMediaFilesForUserThatWereBackedUpInCloud(
	ctx context.Context, requestedFileNames []string, user *models.User,
) (fileKeysThatWereBackedUp []string, fileNamesThatWereBackedUp []string, err error) {
	requestedFileNamesKeys := uc.getMediaFileKeysForUserFromFileNames(user, requestedFileNames)

	userMediaFolderKey := uc.getMediaFolderKeyForUser(user)
	userBackedUpMediaFileKeys, err := uc.MediaCloudService.GetListOfFilesFromFolder(
		ctx, userMediaFolderKey, user.S3MaxNumberOfMediaFiles,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting list of backed up media files for user %s: %v", user.Email, err)
	}

	// sort the userBackedUpMediaFileKeys strings for binary search for faster lookup
	sort.Strings(userBackedUpMediaFileKeys)
	fileKeysThatWereBackedUp = make([]string, 0)
	fileNamesThatWereBackedUp = make([]string, 0)
	for _, requestedFileKey := range requestedFileNamesKeys {
		i := sort.SearchStrings(userBackedUpMediaFileKeys, requestedFileKey)
		if i != len(userBackedUpMediaFileKeys) && userBackedUpMediaFileKeys[i] == requestedFileKey {
			fileKeysThatWereBackedUp = append(fileKeysThatWereBackedUp, requestedFileKey)
			fileNamesThatWereBackedUp = append(fileNamesThatWereBackedUp, requestedFileKey[len(userMediaFolderKey):])
		}
	}

	return fileKeysThatWereBackedUp, fileNamesThatWereBackedUp, nil
}
