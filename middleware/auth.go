package middleware

import (
    "hospital-back/utils"
    "github.com/gofiber/fiber/v2"
    "strings"
)

// Middleware que protege rutas con JWT
func Autenticacion(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{"error": "Falta token de autenticación"})
    }

    partes := strings.Split(authHeader, " ")
    if len(partes) != 2 || partes[0] != "Bearer" {
        return c.Status(401).JSON(fiber.Map{"error": "Formato de token inválido"})
    }

    tokenString := partes[1]

    _, err := utils.ValidarToken(tokenString)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Token inválido o expirado"})
    }

    return c.Next()
}
