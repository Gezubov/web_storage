package repositories

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"log/slog"
	"web_storage/internal/config"
	"web_storage/internal/infrastructure/storage"
	"web_storage/pkg/logger"
)

type MinioRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioRepository(cfgMinio *config.ConfigMinio) *MinioRepository {
	return &MinioRepository{
		client:     storage.MinioClient,
		bucketName: cfgMinio.Bucket,
	}
}

func (repo *MinioRepository) UploadFileMinio(objectName string, fileData *minio.PutObjectOptions, size int64, reader io.Reader) error {
	_, err := repo.client.PutObject(
		context.Background(),
		repo.bucketName,
		objectName,
		reader,
		size,
		*fileData,
	)
	if err != nil {
		logger.Logger.Error("Failed to upload file to MinIO",
			slog.String("bucket", repo.bucketName),
			slog.String("object", objectName),
			slog.String("error", err.Error()))
		return err
	}

	logger.Logger.Info("File uploaded to MinIO successfully",
		slog.String("bucket", repo.bucketName),
		slog.String("object", objectName))

	return nil
}

func (repo *MinioRepository) DownloadFileMinio(objectName string) (*minio.Object, error) {
	object, err := repo.client.GetObject(context.Background(), repo.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Logger.Error("Failed to download file from MinIO",
			slog.String("bucket", repo.bucketName),
			slog.String("object", objectName),
			slog.String("error", err.Error()))
		return nil, err
	}

	logger.Logger.Info("File downloaded from MinIO successfully",
		slog.String("bucket", repo.bucketName),
		slog.String("object", objectName))

	return object, nil
}

func (repo *MinioRepository) DeleteFileMinio(objectName string) error {
	err := repo.client.RemoveObject(context.Background(), repo.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Logger.Error("Failed to delete file from MinIO",
			slog.String("bucket", repo.bucketName),
			slog.String("object", objectName),
			slog.String("error", err.Error()))
		return err
	}
	logger.Logger.Info("File deleted from MinIO successfully",
		slog.String("bucket", repo.bucketName),
		slog.String("object", objectName))
	return nil
}
