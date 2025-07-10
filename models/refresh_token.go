package models

import "time"

// RefreshToken res un token de larga vida para renovar el AccessTokens.

type RefreshToken struct {
    ID        int       `json:"id"`
    Token     string    `json:"token"`
    UserID    int       `json:"user_id"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}
