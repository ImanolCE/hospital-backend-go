package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/pquerna/otp/totp"
    "golang.org/x/crypto/bcrypt"
    "hospital-back/models"
    "hospital-back/utils"
)

type RecoverMfaDTO struct {
    Correo   string `json:"correo"`
    Password string `json:"password"`
}

// RecoverMfaStart valida credenciales y devuelve un nuevo QR/secret
func RecoverMfaStart(c *fiber.Ctx) error {
    var dto RecoverMfaDTO
    if err := c.BodyParser(&dto); err != nil {
        return utils.ResponseError(c, 400, "MFA10", "Datos inválidos")
    }
    // 1) Buscar usuario
    user, err := models.GetUserByEmail(dto.Correo)
    if err != nil {
        return utils.ResponseError(c, 404, "MFA11", "Usuario no encontrado")
    }
    // 2) Verificar contraseña
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)) != nil {
        return utils.ResponseError(c, 401, "MFA12", "Credenciales inválidas")
    }
    // 3) Generar nuevo secreto TOTP
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: user.Correo,
    })
    if err != nil {
        return utils.ResponseError(c, 500, "MFA13", "Error generando OTP")
    }
    secret := key.Secret()
    // 4) Guardar secret y desactivar temporalmente MFA
    if err := models.UpdateUserMfaSecret(user.ID, secret); err != nil {
        return utils.ResponseError(c, 500, "MFA14", "No se pudo actualizar OTP")
    }
    // 5) Devolver la URL para el QR y el secret
    return utils.ResponseSuccess(c, 200, "MFA15", fiber.Map{
        "otpAuthUrl": key.URL(),
        "secret":     secret,
    })
}

