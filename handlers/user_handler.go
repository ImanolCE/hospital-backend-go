
package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "github.com/gofiber/fiber/v2"
)

// CreateUser crrea un nuevo usuario en la bd
func CreateUser(c *fiber.Ctx) error {
    user := new(models.User) 

    // Lee y convierte los datos del body de la petición a la estructura User
    if err := c.BodyParser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // Inserta el usuario en la bd
    _, err := config.DB.Exec(context.Background(),
        "INSERT INTO usuarios (nombre, apellido, correo, contraseña, tipo_usuario) VALUES ($1, $2, $3, $4, $5)",
        user.Nombre, user.Apellido, user.Correo, user.Password, user.Tipo)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al insertar usuario"}) 
    }

    return c.JSON(fiber.Map{"message": "Usuario creado correctamente"})
}
