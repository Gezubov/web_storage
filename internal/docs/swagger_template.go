package docs

import (
	"github.com/gofiber/fiber/v2"
	"html/template"
	"time"
)

const (
	swaggerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script> -->
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script> -->
    <link rel="stylesheet" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
    <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" /> -->
	<style>
		body {
			margin: 0;
		}
	</style>
    <title>Swagger</title>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
          SwaggerUIBundle({
            url: "/docs/swagger.json",
            dom_id: '#swagger-ui',
            presets: [
              SwaggerUIBundle.presets.apis,
              SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
          })
        }
    </script>
</body>
</html>
`
)

func SwaggerUI(c *fiber.Ctx) error {
	// Устанавливаем заголовок Content-Type
	c.Type("html", "utf-8")

	// Создаем шаблон для Swagger UI
	tmpl, err := template.New("swagger").Parse(swaggerTemplate)
	if err != nil {
		// В случае ошибки возвращаем 500
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse template")
	}

	// Выполняем шаблон и выводим результат
	err = tmpl.Execute(c.Response().BodyWriter(), struct {
		Time int64
	}{
		Time: time.Now().Unix(),
	})
	if err != nil {
		// В случае ошибки возвращаем 500
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to render template")
	}

	return nil
}
