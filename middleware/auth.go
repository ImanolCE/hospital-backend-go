package middleware

import (
    "hospital-back/utils"
    "github.com/gofiber/fiber/v2"
    "strings"
)

// Autenticacion recibe los permisos que requiere el endpoint
func Autenticacion(permisosNecesarios []string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        hdr := c.Get("Authorization")
        if hdr == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Token no proporcionado"})
        }
        parts := strings.Split(hdr, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{"error": "Formato de token inválido"})
        }

        claims, err := utils.ValidarAccessToken(parts[1])
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Token inválido o expirado"})
        }

        // Extraer permisos del token
        raw := (*claims)["permisos"].([]interface{})
        userPerms := map[string]bool{}
        for _, p := range raw {
            userPerms[p.(string)] = true
        }
        // Validar cada permiso necesario
        for _, req := range permisosNecesarios {
            if !userPerms[req] {
                return c.Status(403).JSON(fiber.Map{"error": "Permiso denegado"})
            }
        }

        c.Locals("user_id", (*claims)["user_id"])
        return c.Next()
    }
}
