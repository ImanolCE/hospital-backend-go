package handlers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "hospital-back/config"
    "hospital-back/utils"
)

func Refresh(c *fiber.Ctx) error {
    var req struct {
        Token string `json:"refresh_token"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Falta refresh_token"})
    }

    // 1. Buscar en l a BD
    var userID int
    var expiresAt time.Time
    err := config.DB.QueryRow(context.Background(),
        "SELECT user_id, expires_at FROM refresh_tokens WHERE token=$1",
        req.Token,
    ).Scan(&userID, &expiresAt)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "RefreshToken inválido"})
    }

    // 2. Verificar la expiración
    if time.Now().After(expiresAt) {
        return c.Status(401).JSON(fiber.Map{"error": "RefreshToken expirado"})
    }

    // 3. Para volver a cargar permisos del usuario
    var rolNombre string
    config.DB.QueryRow(context.Background(),
        "SELECT tipo_usuario FROM usuarios WHERE id_usuario=$1", userID,
    ).Scan(&rolNombre)

    rows, _ := config.DB.Query(context.Background(),
        `SELECT p.nombre FROM permisos p
         JOIN rol_permisos rp ON p.id = rp.permiso_id
         JOIN roles r        ON r.id = rp.rol_id
         WHERE r.nombre = $1`, rolNombre)
    defer rows.Close()

    permisos := []string{}
    for rows.Next() {
        var p string
        rows.Scan(&p)
        permisos = append(permisos, p)
    }

    // 4. Generar un nuevo Access Token
    newAccessToken, err := utils.GenerarToken(userID, permisos)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al generar AccessToken"})
    }

    

    return c.JSON(fiber.Map{
        "access_token": newAccessToken,
    })
}
