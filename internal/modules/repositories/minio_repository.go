package repositories

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"web_storage/internal/config"
	"web_storage/internal/infrastructure/storage"
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
		log.Printf("Failed to upload file to MinIO: %v", err)
		return err
	}
	return nil
}

func (repo *MinioRepository) DownloadFileMinio(objectName string) (*minio.Object, error) {
	object, err := repo.client.GetObject(context.Background(), repo.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Failed to download file from MinIO: %v", err)
		return nil, err
	}
	return object, nil
}

func (repo *MinioRepository) DeleteFileMinio(objectName string) error {
	err := repo.client.RemoveObject(context.Background(), repo.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Failed to delete file from MinIO: %v", err)
		return err
	}
	return nil
}
