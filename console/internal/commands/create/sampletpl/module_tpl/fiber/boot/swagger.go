package boot

import (
	"cvgo/app"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"os"
	"path/filepath"
)

// Swagger
func SwaggerDoc(fiberApp *fiber.App) {
	env := app.Env()
	port := HttpServerPort
	var domain string
	if env == "" {
		domain = "127.0.0.1"
	} else if env == "alpha" {
		domain = "xxxx"
		port = "80"
	}
	// swagger-doc
	docFilePath := app.Config.GetSwagger().FilePath
	docFileName := filepath.Base(docFilePath)

	fiberApp.Get("/"+docFileName, func(c *fiber.Ctx) error {
		data, _ := os.ReadFile(docFilePath)
		return c.SendString(string(data))
	})
	url := "http://" + domain + ":" + port + "/" + docFileName
	fiberApp.Get("/swagger/*", swagger.New(swagger.Config{URL: url}))
	fmt.Println("  Swagger Doc: http://" + domain + ":" + HttpServerPort + "/swagger")
}
