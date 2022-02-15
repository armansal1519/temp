package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type httpError struct {
	Statuscode int    `json:"statusCode"`
	Error      string `json:"error"`
}

// ErrorHandler is used to catch error thrown inside the routes by ctx.Next(err)
//func ErrorHandler(c *fiber.Ctx, err error) error {
//	// Statuscode defaults to 500
//	code := fiber.StatusInternalServerError
//
//	// Check if it's an fiber.Error type
//	if e, ok := err.(*fiber.Error); ok {
//		code = e.Code
//	}
//
//	return c.Status(code).JSON(&httpError{
//		Statuscode: code,
//		Error:      err.Error(),
//	})
//}

func CustomErrorResponse(statusCode int, devErrCode int, givenErr error, userMsg string, c *fiber.Ctx) error {
	devMsg := fmt.Sprintf("%v", givenErr)
	return c.Status(statusCode).JSON(fiber.Map{
		"devErrCode": devErrCode,
		"devMsg":     devMsg,
		"userMsg":    userMsg,
	})
}

type CError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	DevInfo string `json:"devInfo"`
	UserMsg string `json:"userMsg"`
}
