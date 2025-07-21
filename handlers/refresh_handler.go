package handlers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/xeipuuv/gojsonschema"

    "hospital-back/config"
    "hospital-back/utils"
)

// JSON Schema para validar el refresh_token
const refreshSchema = `{
  "type": "object",
  "required": ["refresh_token"],
  "properties": {
    "refresh_token": { "type": "string", "minLength": 8 }
  },
  "additionalProperties": false
}`

func Refresh(c *fiber.Ctx) error {
    // Validación con JSON Schema
    body := c.Body()
    schemaLoader := gojsonschema.NewStringLoader(refreshSchema)
    docLoader := gojsonschema.NewBytesLoader(body)
    result, err := gojsonschema.Validate(schemaLoader, docLoader)
    if err != nil || !result.Valid() {
        return utils.ResponseError(c, 400, "RF01", "Body inválido o faltan campos")
    }

    var req struct {
        Token string `json:"refresh_token"`
    }
    if err := c.BodyParser(&req); err != nil {
        return utils.ResponseError(c, 400, "RF02", "Falta refresh_token")
    }

    // 1. Buscar en la BD
    var userID int
    var expiresAt time.Time
    err = config.DB.QueryRow(context.Background(),
        "SELECT user_id, expires_at FROM refresh_tokens WHERE token=$1",
        req.Token,
    ).Scan(&userID, &expiresAt)
    if err != nil {
        return utils.ResponseError(c, 401, "RF03", "RefreshToken inválido")
    }

    // 2. Verificar la expiración
    if time.Now().After(expiresAt) {
        return utils.ResponseError(c, 401, "RF04", "RefreshToken expirado")
    }

    // 3. Cargar rol y permisos
    var rolNombre string
    err = config.DB.QueryRow(context.Background(),
        "SELECT tipo_usuario FROM usuarios WHERE id_usuario=$1", userID,
    ).Scan(&rolNombre)
    if err != nil {
        return utils.ResponseError(c, 500, "RF05", "Error obteniendo tipo de usuario")
    }

    rows, err := config.DB.Query(context.Background(),
        `SELECT p.nombre FROM permisos p
         JOIN rol_permisos rp ON p.id = rp.permiso_id
         JOIN roles r ON r.id = rp.rol_id
         WHERE r.nombre = $1`, rolNombre)
    if err != nil {
        return utils.ResponseError(c, 500, "RF06", "Error cargando permisos")
    }
    defer rows.Close()

    permisos := []string{}
    for rows.Next() {
        var p string
        rows.Scan(&p)
        permisos = append(permisos, p)
    }

    // 4. Generar un nuevo Access Token
    accessToken, err := utils.GenerarToken(userID, permisos, rolNombre)
    if err != nil {
        return utils.ResponseError(c, 500, "RF07", "Error generando AccessToken")
    }

    // 5. Rotación del Refresh Token
    _, _ = config.DB.Exec(context.Background(),
        "DELETE FROM refresh_tokens WHERE token=$1", req.Token)

    refreshToken, expiresAt := utils.GenerateRefreshToken()
    _, err = config.DB.Exec(context.Background(),
        `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`,
        refreshToken, userID, expiresAt,
    )
    if err != nil {
        return utils.ResponseError(c, 500, "RF08", "Error guardando nuevo RefreshToken")
    }

    // 6. Devolver ambos tokens
    return utils.ResponseSuccess(c, 200, "RF_OK", []fiber.Map{{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }})
}
