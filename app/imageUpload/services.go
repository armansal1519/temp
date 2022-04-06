package imageUpload

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

func handleFileupload(c *fiber.Ctx) error {

	// parse incomming image file

	file, err := c.FormFile("image")

	if err != nil {
		log.Println("image upload error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})

	}

	// generate new uuid for image name
	uniqueId := uuid.New()

	// remove "- from imageName"

	filename := strings.Replace(uniqueId.String(), "-", "", -1)

	// extract image extension from original file filename

	fileExt := strings.Split(file.Filename, ".")[1]

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)

	// save image to ./images dir
	err = c.SaveFile(file, fmt.Sprintf("./images/%s", image))

	if err != nil {
		log.Println("image save error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	// generate image url to serve to client using CDN

	imageUrl := fmt.Sprintf("https://data.bamachoub.com/api/v1/images/%s", image)

	// create meta data and send to client

	//data := map[string]interface{}{
	//
	//	"imageName": image,
	//	"imageUrl":  imageUrl,
	//	"header":    file.Header,
	//	"size":      file.Size,
	//}

	return c.JSON(fiber.Map{"status": 201, "message": "Image uploaded successfully", "urls": []string{imageUrl}})
}

func handleFileuploadWithName(c *fiber.Ctx) error {

	// parse incomming image file

	file, err := c.FormFile("image")

	if err != nil {
		log.Println("image upload error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})

	}
	f := file.Filename

	// generate new uuid for image name
	//uniqueId := uuid.New()
	//
	//// remove "- from imageName"
	//
	//filename := strings.Replace(uniqueId.String(), "-", "", -1)
	//
	//// extract image extension from original file filename
	//
	//fileExt := strings.Split(file.Filename, ".")[1]
	//
	//// generate image from filename and extension
	//image := fmt.Sprintf("%s.%s", filename, fileExt)

	// save image to ./images dir
	err = c.SaveFile(file, fmt.Sprintf("./images/%s", f))

	if err != nil {
		log.Println("image save error --> ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	// generate image url to serve to client using CDN

	imageUrl := fmt.Sprintf("https://data.bamachoub.com/api/v1/images/%s", f)

	// create meta data and send to client

	//data := map[string]interface{}{
	//
	//	"imageName": image,
	//	"imageUrl":  imageUrl,
	//	"header":    file.Header,
	//	"size":      file.Size,
	//}

	return c.JSON(fiber.Map{"status": 201, "message": "Image uploaded successfully", "urls": []string{imageUrl}})
}

func multiple(c *fiber.Ctx) error {
	imageUrls := make([]string, 0)
	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		if token := form.Value["token"]; len(token) > 0 {
			// Get key value:
			fmt.Println(token[0])
		}

		// Get all files from "documents" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			uniqueId := uuid.New()

			// remove "- from imageName"

			filename := strings.Replace(uniqueId.String(), "-", "", -1)

			// extract image extension from original file filename

			fileExt := strings.Split(file.Filename, ".")[1]

			// generate image from filename and extension
			image := fmt.Sprintf("%s.%s", filename, fileExt)
			imageUrl := fmt.Sprintf("https://data.bamachoub.com/api/v1/images/%s", image)

			imageUrls = append(imageUrls, imageUrl)
			// save image to ./images dir
			err = c.SaveFile(file, fmt.Sprintf("./images/%s", image))
			if err != nil {
				return c.JSON(err)
			}
		}
	}

	return c.JSON(fiber.Map{"status": 201, "message": "Image uploaded successfully", "urls": imageUrls})
}

func handleDeleteImage(c *fiber.Ctx) error {
	// extract image name from params
	imageName := c.Params("imageName")
	// delete image from ./images
	log.Println(imageName)
	imageNameArr := strings.Split(imageName, ",")
	log.Println(imageNameArr)
	errArr := make([]error, 0)
	for _, img := range imageNameArr {
		err := os.Remove(fmt.Sprintf("./images/%s", img))
		if err != nil {
			log.Println(err)
			errArr = append(errArr, err)
		}
		return c.JSON(fiber.Map{"status": 400, "message": "Server Error", "data": errArr})
	}

	return c.JSON(fiber.Map{"status": 201, "message": "Image(s) deleted successfully", "data": nil})
}
