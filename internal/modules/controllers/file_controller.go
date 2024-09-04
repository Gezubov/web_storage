package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log/slog"
	"mime/multipart"
	"web_storage/pkg/logger"

	"github.com/minio/minio-go/v7"
	"web_storage/internal/models"
)

type IFileService interface {
	CreateFileServ(file *multipart.FileHeader) (*models.FileMeta, error)
	GetAllFilesServ() ([]*models.FileMeta, error)
	DownloadFileService(id int) (*models.FileMeta, *minio.Object, error)
	DeleteFileService(id int) error
}

type FileController struct {
	fileService IFileService
}

func NewFileController(fs IFileService) *FileController {
	return &FileController{
		fileService: fs,
	}
}

func (fc *FileController) UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("Failed to get file from form", slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to give your file",
		})
	}

	fileMeta, err := fc.fileService.CreateFileServ(file)
	if err != nil {
		logger.Logger.Error("Failed to upload file", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	logger.Logger.Info("File uploaded successfully", slog.String("filename", fileMeta.Name))
	return c.Status(fiber.StatusCreated).JSON(fileMeta)

}

func (fc *FileController) GetAllFilesHandler(c *fiber.Ctx) error {
	files, err := fc.fileService.GetAllFilesServ()
	if err != nil {
		logger.Logger.Error("Failed to retrieve files", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve files",
		})
	}

	logger.Logger.Info("Files retrieved successfully", slog.Int("file_count", len(files)))
	return c.JSON(files)
}

func (fc *FileController) DownloadFileHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.Logger.Error("Invalid file ID", slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}
	// Используем сервис для получения файла и метаданных
	fileMeta, object, err := fc.fileService.DownloadFileService(id)
	if err != nil {
		logger.Logger.Error("Failed to download file", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to download file",
		})
	}
	if fileMeta == nil {
		logger.Logger.Warn("File not found", slog.Int("file_id", id))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		logger.Logger.Error("Failed to get file information from storage", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get file information from storage",
		})
	}

	// Устанавливаем заголовки для скачивания файла с оригинальным именем
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.Name))
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Length", fmt.Sprintf("%d", stat.Size))

	// Передаем файл пользователю
	if _, err := io.Copy(c.Response().BodyWriter(), object); err != nil {
		logger.Logger.Error("Failed to send file to client", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send file to client",
		})
	}

	logger.Logger.Info("File downloaded successfully", slog.String("filename", fileMeta.Name))
	return nil
}

func (fc *FileController) DeleteFileHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.Logger.Error("Invalid file ID", slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}
	// Удаляем запись о файле из базы данных
	if err := fc.fileService.DeleteFileService(id); err != nil {
		logger.Logger.Error("Failed to delete file", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file from database",
		})
	}

	logger.Logger.Info("File deleted successfully", slog.Int("file_id", id))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
