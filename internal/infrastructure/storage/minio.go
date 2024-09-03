package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"web_storage/internal/config"
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
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Убедитесь, что существует бакет для хранения файлов

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Проверка, если бакет уже существует
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	MinioClient = minioClient
}
