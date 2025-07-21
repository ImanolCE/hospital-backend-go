package models

import (
    "context"
    "hospital-back/config"
)

// User representa un registro de la tabla usuarios
type User struct {
    ID         int    `json:"id_usuario"`
    Nombre     string `json:"nombre"`
    Apellido   string `json:"apellido"`
    Correo     string `json:"correo"`
    Password   string `json:"password"`
    Tipo       string `json:"tipo_usuario"`
    MFASecret  string `json:"mfa_secret"`
    MFAEnabled bool   `json:"mfa_enabled"`
}

// GetUserByID obtiene un usuario por su ID, incluyendo MFA fields
func GetUserByID(userID int) (*User, error) {
    var u User
    err := config.DB.QueryRow(
        context.Background(),
        `SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario, mfa_secret, mfa_enabled
         FROM usuarios
         WHERE id_usuario = $1`,
        userID,
    ).Scan(
        &u.ID,
        &u.Nombre,
        &u.Apellido,
        &u.Correo,
        &u.Password,
        &u.Tipo,
        &u.MFASecret,
        &u.MFAEnabled,
    )
    if err != nil {
        return nil, err
    }
    return &u, nil
}

// GetUserByEmail busca un usuario por correo, incluyendo MFA fields
func GetUserByEmail(email string) (*User, error) {
    var u User
    err := config.DB.QueryRow(
        context.Background(),
        `SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario, mfa_secret, mfa_enabled
         FROM usuarios
         WHERE correo = $1`,
        email,
    ).Scan(
        &u.ID,
        &u.Nombre,
        &u.Apellido,
        &u.Correo,
        &u.Password,
        &u.Tipo,
        &u.MFASecret,
        &u.MFAEnabled,
    )
    if err != nil {
        return nil, err
    }
    return &u, nil
}

// UpdateUserMfaSecret actualiza el secret TOTP y desactiva MFA temporalmente
func UpdateUserMfaSecret(userID int, secret string) error {
    _, err := config.DB.Exec(
        context.Background(),
        `UPDATE usuarios SET mfa_secret = $1, mfa_enabled = false WHERE id_usuario = $2`,
        secret, userID,
    )
    return err
}

// ActivateUserMfa marca MFA como activado para el usuario
func ActivateUserMfa(userID int) error {
    _, err := config.DB.Exec(
        context.Background(),
        `UPDATE usuarios SET mfa_enabled = true WHERE id_usuario = $1`,
        userID,
    )
    return err
}
