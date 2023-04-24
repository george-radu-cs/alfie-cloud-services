package services

import (
	"api/app/utils"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"
)

type mediaCloudService struct {
	s3BucketName               *string
	s3Client                   *s3.Client
	s3PresignClient            *s3.PresignClient
	putFileExpirationTime      time.Duration
	downloadFileExpirationTime time.Duration
}

func NewMediaCloudService() MediaCloudService {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		utils.ErrorLogger.Fatalf("unable to load aws s3 SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)

	putFileExpirationTime, err := time.ParseDuration(os.Getenv("PUT_FILE_EXPIRATION_TIME"))
	if err != nil {
		utils.ErrorLogger.Fatalf("unable to parse env var PUT_FILE_EXPIRATION_TIME as time.Duration: %v", err)
	}
	downloadFileExpirationTime, err := time.ParseDuration(os.Getenv("DOWNLOAD_FILE_EXPIRATION_TIME"))
	if err != nil {
		utils.ErrorLogger.Fatalf("unable to parse env var DOWNLOAD_FILE_EXPIRATION_TIME as time.Duration: %v", err)
	}

	return &mediaCloudService{
		s3BucketName:               aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		s3Client:                   s3Client,
		s3PresignClient:            s3PresignClient,
		putFileExpirationTime:      putFileExpirationTime,
		downloadFileExpirationTime: downloadFileExpirationTime,
	}
}

func (mcs *mediaCloudService) CreatePresignedURLForFileUpload(fileName string) (presignedUploadURL string, err error) {
	request, err := mcs.s3PresignClient.PresignPutObject(
		context.TODO(), &s3.PutObjectInput{
			Bucket: mcs.s3BucketName,
			Key:    aws.String(fileName),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = mcs.putFileExpirationTime
		},
	)
	if err != nil {
		return "", fmt.Errorf("couldn't create a presigned request to put %s: %v", fileName, err)
	}

	return request.URL, nil
}

func (mcs *mediaCloudService) CreatePresignedURLForFileDownload(fileName string) (
	presignedDownloadURL string, err error,
) {
	request, err := mcs.s3PresignClient.PresignGetObject(
		context.TODO(), &s3.GetObjectInput{
			Bucket: mcs.s3BucketName,
			Key:    aws.String(fileName),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = mcs.downloadFileExpirationTime
		},
	)
	if err != nil {
		return "", fmt.Errorf("couldn't create a presigned request to get %s: %v", fileName, err)
	}

	return request.URL, nil
}

func (mcs *mediaCloudService) DeleteFile(fileName string) (err error) {
	_, err = mcs.s3Client.DeleteObject(
		context.TODO(), &s3.DeleteObjectInput{
			Bucket: mcs.s3BucketName,
			Key:    aws.String(fileName),
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't delete %s: %v", fileName, err)
	}

	return nil
}

func (mcs *mediaCloudService) CreatePresignedURLsForMultipleFilesUpload(fileNames []string) (
	presignedUploadURLs []string, err error,
) {
	presignedUploadURLs = make([]string, len(fileNames))
	for i, fileName := range fileNames {
		presignedUploadLink, err := mcs.CreatePresignedURLForFileUpload(fileName)
		if err != nil {
			return nil, err
		}

		presignedUploadURLs[i] = presignedUploadLink
	}

	return presignedUploadURLs, nil
}

func (mcs *mediaCloudService) CreatePresignedURLsForMultipleFilesDownload(fileNames []string) (
	presignedDownloadURLs []string, err error,
) {
	presignedDownloadURLs = make([]string, len(fileNames))
	for i, fileName := range fileNames {
		presignedDownloadLink, err := mcs.CreatePresignedURLForFileDownload(fileName)
		if err != nil {
			return nil, err
		}

		presignedDownloadURLs[i] = presignedDownloadLink
	}

	return presignedDownloadURLs, nil
}

func (mcs *mediaCloudService) DeleteMultipleFiles(fileNames []string) (err error) {
	for _, fileName := range fileNames {
		err := mcs.DeleteFile(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mcs *mediaCloudService) CheckIfFileExists(fileNameKey string) (fileExists bool, err error) {
	_, err = mcs.s3Client.HeadObject(
		context.TODO(), &s3.HeadObjectInput{
			Bucket: mcs.s3BucketName,
			Key:    aws.String(fileNameKey),
		},
	)
	if err != nil {
		return false, fmt.Errorf("couldn't find %s: %v", fileNameKey, err)
	}

	return true, nil
}

func (mcs *mediaCloudService) CreateFolder(folderName string) (err error) {
	// in s3 folders don't really exists in the same way as in a file system they are just objects with a slash at the
	// end of their name which allows us to create a hierarchy of objects, since we can get a list of objects by
	// querying from a prefix key

	folderNameKey := mcs.getFolderNameKey(folderName)
	_, err = mcs.s3Client.PutObject(
		context.TODO(), &s3.PutObjectInput{
			Bucket: mcs.s3BucketName,
			Key:    aws.String(folderNameKey),
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create folder %s: %v", folderNameKey, err)
	}

	return nil
}

func (mcs *mediaCloudService) GetListOfFilesFromFolder(folderName string, maxNumberOfFiles int32) (
	fileKeys []string, err error,
) {
	// in s3 folders don't really exists in the same way as in a file system they are just objects with a slash at the
	// end of their name which allows us to create a hierarchy of objects, since we can get a list of objects by
	// querying from a prefix key

	folderNameKey := mcs.getFolderNameKey(folderName)
	response, err := mcs.s3Client.ListObjectsV2(
		context.TODO(), &s3.ListObjectsV2Input{
			Bucket:     mcs.s3BucketName,
			StartAfter: aws.String(folderNameKey),
			MaxKeys:    maxNumberOfFiles,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't get list of files in folder %s: %v", folderNameKey, err)
	}

	fileKeys = make([]string, response.KeyCount)
	for i, object := range response.Contents {
		fileKeys[i] = *object.Key
	}

	return fileKeys, nil
}

func (mcs *mediaCloudService) getFolderNameKey(folderName string) (folderKey string) {
	// make sure the folder name ends with a slash, otherwise the object won't be treated as a folder, but an empty file
	if folderName[len(folderName)-1] != '/' {
		return folderName + "/"
	}
	return folderName
}
