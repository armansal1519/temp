package main

import (
	"bamachoub-backend-go-v1/config"
	_ "bamachoub-backend-go-v1/docs"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

// @title Bamachoub Application
// @version 2.0
// @description This is an API for Bamachoub Application
// @contact.name Arman Salehi
// @contact.email armansal1519@gmail.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func main() {
	app := fiber.New(fiber.Config{
		//Prefork: true,
	})
	app.Use(cors.New())
	//app.Use(etag.New())
	app.Use(pprof.New())
	app.Use(logger.New())
	app.Use(recover.New())

	//app.Use(cache.New(cache.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.Query("refresh") == "true"
	//	},
	//	Expiration: 30 * time.Minute,
	//	CacheControl: true,
	//}))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/swagger/*", swagger.Handler) // default

	InitRoutes(v1)

	log.Fatal(app.Listen(fmt.Sprintf(":%v", config.PORT)))

}
