package imageUpload

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.Router) {
	r := app.Group("/images")

	r.Static("", "./images")

	// handle image uploading using post request

	r.Post("", handleFileupload)
	r.Post("/name", handleFileuploadWithName)
	r.Post("/m", multiple)

	// delete uploaded image by providing unique image name

	r.Delete("/:imageName", handleDeleteImage)

}
