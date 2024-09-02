package router

import (
	"github.com/gofiber/fiber/v2"
	"web_storage/internal/modules/controllers"
)

const (
	files   = "files"
	filesID = "files/:id"
)

func SetupRouter(app *fiber.App, fileController *controllers.FileController) {
	app.Get(filesID, fileController.DownloadFileHandler)
	app.Delete(filesID, fileController.DeleteFileHandler)
	app.Get(files, fileController.GetAllFilesHandler)
	app.Post(files, fileController.UploadFileHandler)
}
