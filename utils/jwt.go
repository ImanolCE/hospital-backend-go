package utils

import (
    "time"
    "errors"
    "github.com/golang-jwt/jwt/v4"
)

// SecretKey es la clave secreta para firmar los tokens 
var SecretKey = []byte("claveSuperSecreta")

// se genera un token JWT que contiene el ID del usuario
func GenerarToken(userID int) (string, error) {
    // Claims -> los datos que irán dentro del token
    claims := jwt.MapClaims{
        "user_id": userID,                             // se guard el ID del usuario
        "exp":     time.Now().Add(time.Hour * 1).Unix(), // se expira en 1 hora
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(SecretKey)
}

// ValidarToken recibe el token y lo parsea para que devuelva los claims si es válido
func ValidarToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return SecretKey, nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("token inválido o expirado")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("no se pudieron obtener los claims")
    }

    return claims, nil
}
    