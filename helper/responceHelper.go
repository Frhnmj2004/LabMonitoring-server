package helper

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(data interface{}) fiber.Map {
	return fiber.Map{
		"status": 200,
		"data":   data,
	}
}

func ErrorResponse(errorMsg string) fiber.Map {
	return fiber.Map{
		"status": 400,
		"error":  errorMsg,
	}
}
