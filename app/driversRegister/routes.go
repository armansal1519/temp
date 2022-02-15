package driversRegister

import (
	"bamachoub-backend-go-v1/utils"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/car-category")

	r.Post("/", createCarTypeCategory)
	r.Delete("/:key", deleteCarCategory)
	r.Get("/", getAllCarTypes)
}

func DriversRoutes(app fiber.Router) {
	r := app.Group("drivers")

	r.Post("/", func(c *fiber.Ctx) error {
		newDriver := new(driverInfo)
		if err := utils.ParseBodyAndValidate(c, newDriver); err != nil {
			return c.JSON(err)
		}
		resp, err := CreateDriversInfo(*newDriver)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(resp)
	})
	r.Get("/", getAllDriversInfo)
	r.Get("/:key", getDriversInfoByKey)

	r.Delete("/:key", deleteDriversInfoByKey)

	r.Put("/:key", updateDriversInfoByKey)

	r.Post("/search", searchIntoDrivers)
	r.Post("/filter", filterDriversBaseOnCarType)

}
