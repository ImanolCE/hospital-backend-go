package handlers

import (
    "context"
    "net/http"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/pquerna/otp/totp"
    "github.com/xeipuuv/gojsonschema"

    "hospital-back/config"
    "hospital-back/utils"
    "hospital-back/models"
    
)

// JSON Schema para verificar OTP
const verifyMFASchema = `{
  "type": "object",
  "required": ["otp"],
  "properties": {
    "otp": { "type": "string", "minLength": 6, "maxLength": 6 }
  },
  "additionalProperties": false
}`

// MFASetup genera una clave secreta y otpauth_url
func MFASetup(c *fiber.Ctx) error {
    hdr := c.Get("Authorization")
    parts := strings.Split(hdr, " ")
    claims, err := utils.ValidarAccessToken(parts[1])
    if err != nil {
        return utils.ResponseError(c, http.StatusUnauthorized, "MFA01", "Token inválido")
    }
    userID := int((*claims)["user_id"].(float64))

    // Obtener correo
    var email string
    if err := config.DB.QueryRow(context.Background(),
        "SELECT correo FROM usuarios WHERE id_usuario=$1", userID).
        Scan(&email); err != nil {
        return utils.ResponseError(c, 500, "MFA02", "Usuario no encontrado")
    }

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: email,
    })
    if err != nil {
        return utils.ResponseError(c, 500, "MFA03", "Error generando MFA")
    }

    // Guardar secret (mfa_enabled sigue como estaba)
    if _, err := config.DB.Exec(context.Background(),
        "UPDATE usuarios SET mfa_secret=$1 WHERE id_usuario=$2",
        key.Secret(), userID,
    ); err != nil {
        return utils.ResponseError(c, 500, "MFA04", "Error guardando secret")
    }

    return utils.ResponseSuccess(c, 200, "MFA_SETUP", []fiber.Map{{
        "mfa_secret":  key.Secret(),
        "otpauth_url": key.URL(),
    }})
}

// MFAVerify verifica un código OTP y activa MFA si es correcto
func MFAVerify(c *fiber.Ctx) error {
    hdr := c.Get("Authorization")
    parts := strings.Split(hdr, " ")
    claims, err := utils.ValidarAccessToken(parts[1])
    if err != nil {
        return utils.ResponseError(c, 401, "MFA01", "Token inválido")
    }
    userID := int((*claims)["user_id"].(float64))

    body := c.Body()
    schemaLoader := gojsonschema.NewStringLoader(verifyMFASchema)
    docLoader := gojsonschema.NewBytesLoader(body)
    result, err := gojsonschema.Validate(schemaLoader, docLoader)
    if err != nil {
        return utils.ResponseError(c, 400, "MFA06", "Error validando OTP")
    }
    if !result.Valid() {
        return utils.ResponseError(c, 400, "MFA07", "OTP inválido (schema)")
    }

    var input struct {
        OTP string `json:"otp"`
    }
    if err := c.BodyParser(&input); err != nil || input.OTP == "" {
        return utils.ResponseError(c, 400, "MFA08", "OTP requerido")
    }

    var secret string
    config.DB.QueryRow(context.Background(),
        "SELECT mfa_secret FROM usuarios WHERE id_usuario=$1", userID,
    ).Scan(&secret)

    if !totp.Validate(input.OTP, secret) {
        return utils.ResponseError(c, 401, "MFA09", "OTP incorrecto")
    }

    if _, err := config.DB.Exec(context.Background(),
        "UPDATE usuarios SET mfa_enabled=true WHERE id_usuario=$1", userID,
    ); err != nil {
        return utils.ResponseError(c, 500, "MFA10", "Error activando MFA")
    }

    return utils.ResponseSuccess(c, 200, "MFA_VERIFIED", []fiber.Map{{
        "message": "MFA activado correctamente",
    }})
}

// RegenerateMfa genera un nuevo secret y QR para el usuario autenticado
func RegenerateMfa(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)

    // 1) Busca el email del usuario (para el label del QR)
    user, err := models.GetUserByID(userID)
    if err != nil {
        return utils.ResponseError(c, 404, "MFA04", "Usuario no encontrado")
    }

    // 2) Genera un nuevo secreto TOTP
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: user.Correo,
    })
    if err != nil {
        return utils.ResponseError(c, 500, "MFA05", "Error generando OTP")
    }

    secret := key.Secret()

    // 3) Guarda el nuevo secreto en BD
    if err := models.UpdateUserMfaSecret(userID, secret); err != nil {
        return utils.ResponseError(c, 500, "MFA06", "No se pudo actualizar OTP")
    }

    // 4) Devuelve la URI para el QR (puedes convertirla a imagen en el frontend)
    return utils.ResponseSuccess(c, 200, "MFA03", fiber.Map{
        "otpAuthUrl": key.URL(),
        "secret":     secret, // opcional si quieres mostrar el código manual
    })
}

func ActivarMFA(c *fiber.Ctx) error {
  var body struct{ Otp string }
  if err := c.BodyParser(&body); err != nil {
    return utils.ResponseError(c,400,"MFA20","Datos inválidos")
  }
  userID := c.Locals("user_id").(int)
  user, err := models.GetUserByID(userID)
  if err != nil {
    return utils.ResponseError(c,404,"MFA21","Usuario no encontrado")
  }
  if user.MFASecret == "" {
    return utils.ResponseError(c,400,"MFA22","MFA no iniciada")
  }
  if !totp.Validate(body.Otp, user.MFASecret) {
    return utils.ResponseError(c,401,"MFA23","OTP inválido")
  }
  if err := models.ActivateUserMfa(userID); err != nil {
    return utils.ResponseError(c,500,"MFA24","No se pudo activar MFA")
  }
  return utils.ResponseSuccess(c,200,"MFA25",nil)
}


// ActivateMfaRecovery activa MFA recien creada por /mfa/recovery
func ActivateMfaRecovery(c *fiber.Ctx) error {
  var dto struct {
    Correo string `json:"correo"`
    OTP    string `json:"otp"`
  }
  if err := c.BodyParser(&dto); err != nil {
    return utils.ResponseError(c, 400, "MFA03", "Body inválido o faltan campos")
  }
  // 1) Leer secret
  var secret string
  err := config.DB.QueryRow(context.Background(),
    "SELECT mfa_secret FROM usuarios WHERE correo=$1", dto.Correo).
    Scan(&secret)
  if err != nil || secret == "" {
    return utils.ResponseError(c, 404, "MFA04", "Usuario no encontrado o sin MFA")
  }
  // 2) Validar OTP
  if !totp.Validate(dto.OTP, secret) {
    return utils.ResponseError(c, 401, "MFA05", "OTP inválido")
  }
  // 3) Marcar MFA habilitado
  if _, err := config.DB.Exec(context.Background(),
    "UPDATE usuarios SET mfa_enabled = true WHERE correo = $1", dto.Correo); err != nil {
    return utils.ResponseError(c, 500, "MFA06", "Error al activar MFA")
  }
  return utils.ResponseSuccess(c, 200, "MFA07", fiber.Map{
    "message": "MFA activado correctamente",
  })
}

