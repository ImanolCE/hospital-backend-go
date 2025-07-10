package utils

import "github.com/gofiber/fiber/v2"

// ResponseSuccess envuelve la respuesta exitosa y el data puede ser cualquier cosa (array u objeto).
func ResponseSuccess(c *fiber.Ctx, statusCode int, intCode string, data interface{}) error {
    return c.Status(statusCode).JSON(fiber.Map{
        "statusCode": statusCode,
        "intCode":    intCode,
        "data":       data,
    })
}

// ResponseError y envuelve los errores de la API.
func ResponseError(c *fiber.Ctx, statusCode int, intCode, message string) error {
    return c.Status(statusCode).JSON(fiber.Map{
        "statusCode": statusCode,
        "intCode":    intCode,
        "error":      message,
    })
}
