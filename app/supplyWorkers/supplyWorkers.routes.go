package supplyWorkers

import "github.com/gofiber/fiber/v2"

func Routes(app fiber.Router) {
	r := app.Group("/supplier-employee")
	//r.Get("", GetWarehouse)
	r.Post("/manager", CreateSupplyManager)
	r.Get("/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		return c.JSON(GetSupplyWorkerByKey(key))
	})

}
