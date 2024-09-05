package router

import (
	"github.com/gofiber/fiber/v2"
	"web_storage/internal/docs"
	_ "web_storage/internal/docs"
	"web_storage/internal/modules/controllers"
)

const (
	files   = "files"
	filesID = "files/:id"
)

func SetupRouter(app *fiber.App, fileController *controllers.FileController) {
	app.Get("/swagger/*", docs.SwaggerUI)
	app.Static("/docs", "./") // Сервис swagger.json как статического файла
	app.Get("/docs/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	app.Get(filesID, fileController.DownloadFileHandler)
	app.Delete(filesID, fileController.DeleteFileHandler)
	app.Get(files, fileController.GetAllFilesHandler)
	app.Post(files, fileController.UploadFileHandler)
}
