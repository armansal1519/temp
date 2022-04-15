package similarityGraph

import (
	"bamachoub-backend-go-v1/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/p-similarity")

	r.Post("/node", middleware.Auth, createSimilarityNode)
	r.Post("/node/admin", middleware.CheckAdmin, createSimilarityNodeByAdmin)
	r.Post("/edge", createSimilarityEdge)
	r.Post("/:op/:productKey/:nodeKey", AddOrRemoveProductToNode)
	r.Get("/", getAllSimilarityNodes)
	r.Get("/node/:nodeKey", getOneNodeByKey)
	r.Get("/graph-sim/:key", getSimilarNodeToOneNodeByNodeKey)
	r.Get("/near-nodes/:key", getNearNodesWithProductKey)

	r.Get("/:key", func(c *fiber.Ctx) error {
		offset := c.Query("offset")
		limit := c.Query("limit")
		key := c.Params("key")
		resp, err := getSimilarProductsByProductKey(key, offset, limit)
		if err != nil {
			return c.JSON(err)

		}
		return c.JSON(resp)

	})
	r.Put("/:key", middleware.CheckAdmin, middleware.Auth, updateSimilarityNode)
	r.Delete("/:key", middleware.CheckAdmin, middleware.Auth, removeNode)

}
