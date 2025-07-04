package utils

import "github.com/gofiber/fiber/v2"

type ResponseTemplate struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseErrorTemplate struct {
	Error string `json:"error"`
}

func ResponseSuccess(code int, message string, data interface{}) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(code).JSON(ResponseTemplate{
			Code:    code,
			Message: message,
			Data:    data,
		})
	}
}

func ResponseError(code int, message string, err error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(code).JSON(ResponseTemplate{
			Code:    code,
			Message: message,
			Data: ResponseErrorTemplate{
				Error: err.Error(),
			},
		})
	}
}
