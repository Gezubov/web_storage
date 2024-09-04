package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"web_storage/internal/config"
	"web_storage/internal/infrastructure/db/postgres"
	"web_storage/internal/infrastructure/storage"
	"web_storage/internal/modules/controllers"
	"web_storage/internal/modules/repositories"
	"web_storage/internal/modules/services"
	"web_storage/internal/router"
	"web_storage/pkg/logger"
)

func main() {

	logger.InitLogger()
	logger.Logger.Info("Logger initialized")

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Logger.Error("Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	storage.InitMinio(&cfg.Minio) // Инициализация MinIO

	database, err := postgres.Connect(cfg)
	if err != nil {
		logger.Logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Создаем репозиторий сервис и контроллер
	minioRepository := repositories.NewMinioRepository(&cfg.Minio)
	fileRepository := repositories.NewFileRepository(database)
	fileService := services.NewFileService(fileRepository, minioRepository)
	fileController := controllers.NewFileController(fileService)

	logger.Logger.Info("File service and controller initialized")

	app := fiber.New()
	app.Use(cors.New())

	router.SetupRouter(app, fileController)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			logger.Logger.Error("Server shutdown with error", slog.String("error", err.Error()))
		}
	}()

	logger.Logger.Info("Server is running", slog.String("port", cfg.AppPort))
	<-quit
	logger.Logger.Info("Shutting down the server")

	// Закрываем сервер
	if err := app.Shutdown(); err != nil {
		logger.Logger.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	// Закрываем соединение с базой данных
	if err := database.Close(); err != nil {
		logger.Logger.Error("Failed to close database connection", slog.String("error", err.Error()))
	}

	logger.Logger.Info("Server exited properly")

}
