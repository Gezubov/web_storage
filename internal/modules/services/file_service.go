package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"
	"web_storage/internal/models"
)

type IFileRepository interface {
	CreateFileRepo(file *models.FileMeta) error
	GetAllFilesRepo() ([]*models.FileMeta, error)
	GetFileByIdRepo(id int) (*models.FileMeta, error)
	DeleteFileByIdRepo(id int) error
}

type IMinioRepository interface {
	UploadFileMinio(objectName string, fileData *minio.PutObjectOptions, size int64, reader io.Reader) error
	DeleteFileMinio(objectName string) error
	DownloadFileMinio(objectName string) (*minio.Object, error)
}

type FileService struct {
	fileRepository  IFileRepository
	minioRepository IMinioRepository
}

func NewFileService(fr IFileRepository, mr IMinioRepository) *FileService {
	return &FileService{
		fileRepository:  fr,
		minioRepository: mr,
	}
}

func (fs *FileService) CreateFileServ(file *multipart.FileHeader) (*models.FileMeta, error) {
	// Хеширование имени файла
	hash := sha256.New()
	hash.Write([]byte(file.Filename + time.Now().String()))
	hashedFileName := hex.EncodeToString(hash.Sum(nil)) + filepath.Ext(file.Filename)

	// Открытие файла для чтения
	fileData, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileData.Close()

	// Загрузка файла в MinIO
	bucketName := "uploads"
	objectName := hashedFileName
	fileSize := file.Size

	err = fs.minioRepository.UploadFileMinio(objectName, &minio.PutObjectOptions{ContentType: "application/octet-stream"}, fileSize, fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Создание записи о файле в базе данных
	fileMeta := &models.FileMeta{
		Name: file.Filename,
		Size: file.Size,
		Link: fmt.Sprintf("/%s/%s", bucketName, objectName),
	}

	if err := fs.fileRepository.CreateFileRepo(fileMeta); err != nil {
		return nil, fmt.Errorf("failed to save file information: %w", err)
	}

	return fileMeta, nil
}

func (fs *FileService) GetAllFilesServ() ([]*models.FileMeta, error) {
	return fs.fileRepository.GetAllFilesRepo()
}

func (fs *FileService) DownloadFileService(id int) (*models.FileMeta, *minio.Object, error) {
	// Получаем информацию о файле
	fileMeta, err := fs.fileRepository.GetFileByIdRepo(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file information: %w", err)
	}
	if fileMeta == nil {
		return nil, nil, nil
	}

	// Загружаем файл из MinIO
	objectName := filepath.Base(fileMeta.Link)
	object, err := fs.minioRepository.DownloadFileMinio(objectName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file from storage: %w", err)
	}

	return fileMeta, object, nil
}

func (fs *FileService) DeleteFileService(id int) error {
	// Получаем информацию о файле
	fileMeta, err := fs.fileRepository.GetFileByIdRepo(id)
	if err != nil {
		return fmt.Errorf("failed to get file information: %w", err)
	}
	if fileMeta == nil {
		return fmt.Errorf("file not found")
	}

	// Удаляем файл из MinIO
	objectName := filepath.Base(fileMeta.Link)
	err = fs.minioRepository.DeleteFileMinio(objectName)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Удаляем запись о файле из базы данных
	err = fs.fileRepository.DeleteFileByIdRepo(id)
	if err != nil {
		return fmt.Errorf("failed to delete file from database: %w", err)
	}

	return nil
}
