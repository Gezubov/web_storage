package docs

import (
	"mime/multipart"
	"web_storage/internal/models"
)

// swagger:route POST /files files uploadFileRequest
// Загрузка файла.
// responses:
//   201: uploadFileResponse
//   400: badRequestResponse
//   500: internalServerErrorResponse

// swagger:parameters uploadFileRequest
type uploadFileRequest struct {
	// in:formData
	// description: Файл для загрузки.
	// required: true
	// swagger:file
	File multipart.FileHeader `json:"file"`
}

// swagger:response uploadFileResponse
type uploadFileResponse struct {
	// in:body
	Body models.FileMeta
}

// swagger:route GET /files files getAllFilesRequest
// Получение всех загруженных файлов.
// responses:
//   200: getAllFilesResponse
//   500: internalServerErrorResponse

// swagger:response getAllFilesResponse
type getAllFilesResponse struct {
	// in:body
	Body []models.FileMeta
}

// swagger:route GET /files/{id} files downloadFileRequest
// Скачивание файла по ID.
// responses:
//   200: downloadFileResponse
//   404: notFoundResponse
//   500: internalServerErrorResponse

// swagger:parameters downloadFileRequest
type downloadFileRequest struct {
	// in:path
	// description: ID файла.
	// required: true
	// type: integer
	ID int `json:"id"`
}

// swagger:response downloadFileResponse
type downloadFileResponse struct {
	// in:body
	Body models.FileMeta
}

// swagger:route DELETE /files/{id} files deleteFileRequest
// Удаление файла по ID.
// responses:
//   200: deleteFileResponse
//   404: notFoundResponse
//   500: internalServerErrorResponse

// swagger:parameters deleteFileRequest
type deleteFileRequest struct {
	// in:path
	// description: ID файла.
	// required: true
	// type: integer
	ID int `json:"id"`
}

// swagger:response deleteFileResponse
type deleteFileResponse struct {
	// description: Успешное удаление файла.
	// in:body
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:response badRequestResponse
type badRequestResponse struct {
	// Описание ошибки 400
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:response internalServerErrorResponse
type internalServerErrorResponse struct {
	// Описание ошибки 500
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// swagger:response notFoundResponse
type notFoundResponse struct {
	// Описание ошибки 404
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}
