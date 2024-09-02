package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"web_storage/internal/infrastructure/db/postgres"
	"web_storage/internal/infrastructure/storage"
	"web_storage/internal/modules/controllers"
	"web_storage/internal/modules/repositories"
	"web_storage/internal/modules/services"
	"web_storage/internal/router"
)

func main() {
	storage.InitMinio() // Инициализация MinIO
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	appPort := os.Getenv("APP_PORT")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	database, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем репозиторий сервис и контроллер
	minioRepository := repositories.NewMinioRepository(bucketName)
	fileRepository := repositories.NewFileRepository(database)
	fileService := services.NewFileService(fileRepository, minioRepository)
	fileController := controllers.NewFileController(fileService)

	app := fiber.New()
	app.Use(cors.New())

	router.SetupRouter(app, fileController)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + appPort); err != nil {
			log.Printf("Server shutdown with error: %v", err)
		}
	}()

	<-quit

	log.Println("Shutdown Server ...")

	// Закрываем сервер
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	// Закрываем соединение с базой данных
	if err := database.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}

	log.Println("Server exited properly")

}
