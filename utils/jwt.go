package utils

import (
    "time"
    //"errors"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
)

// SecretKey es la clave secreta para firmar los tokens 
var SecretKey = []byte("claveSecretaHospital")


// Genera un token que incluye user_id y lista de permisos
func GenerarToken(userID int, permisos []string, tipoUsuario string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  userID,
        "permisos": permisos,
        "tipo_usuario": tipoUsuario,
        "exp":      time.Now().Add(5 * time.Minute).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(SecretKey)
}

// ValidarToken recibe el token y parsea para que devuelva los claims si es válido
func ValidarAccessToken(tokenString string) (*jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return SecretKey, nil
    })

    if err != nil || !token.Valid {
        return nil, err
    }
    claims := token.Claims.(jwt.MapClaims)
    return &claims, nil
}   

// GenerateRefreshToken crea un UUID como refresh token .
func GenerateRefreshToken() (string, time.Time) {
    token := uuid.NewString()

    // para que caduque en 7 
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    return token, expiresAt
}