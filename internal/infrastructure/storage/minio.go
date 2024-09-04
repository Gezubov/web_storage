package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log/slog"
	"web_storage/internal/config"
	"web_storage/pkg/logger"
)

var MinioClient *minio.Client

func InitMinio(configMinio *config.ConfigMinio) {

	endPoint := configMinio.Endpoint
	accessKeyID := configMinio.AccessKey
	secretAccessKey := configMinio.SecretKey
	bucketName := configMinio.Bucket
	location := configMinio.Location

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Logger.Error("Failed to initialize MinIO client:", slog.String("error", err.Error()))
		return
	}

	// Убедитесь, что существует бакет для хранения файлов

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Проверка, если бакет уже существует
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			logger.Logger.Info(`Bucket already exists`, slog.String("bucket", bucketName))
		} else {
			logger.Logger.Error("Failed to create bucket:", slog.String("error", err.Error()))
			return
		}
	} else {
		logger.Logger.Info(`Successfully created bucket`, slog.String("bucket", bucketName))
	}

	MinioClient = minioClient
	logger.Logger.Info("MinIO client initialized successfully")
}
