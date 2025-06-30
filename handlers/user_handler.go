
package handlers

import (
    "context"
    "hospital-back/config"
    "hospital-back/models"
    "hospital-back/utils"
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
        "INSERT INTO usuarios (nombre, apellido, correo, password, tipo_usuario) VALUES ($1, $2, $3, $4, $5)",
        user.Nombre, user.Apellido, user.Correo, user.Password, user.Tipo)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al insertar usuario"}) 
    }

    return c.JSON(fiber.Map{"message": "Usuario creado correctamente"})
}


// GetUsers devuelve la lista de usuarios registrados en la base de datos
func GetUsers(c *fiber.Ctx) error {
    rows, err := config.DB.Query(context.Background(), 
    "SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario FROM usuarios")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener usuarios"})
    }
    defer rows.Close()

    var usuarios []models.User

    for rows.Next() {
    var user models.User
    err := rows.Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Correo, &user.Password, &user.Tipo)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos"})
    }
    usuarios = append(usuarios, user)
}


    return c.JSON(usuarios)
}


// Obtener usuario por ID
func GetUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User

   err := config.DB.QueryRow(context.Background(),
    "SELECT id_usuario, nombre, apellido, correo, password, tipo_usuario FROM usuarios WHERE id_usuario=$1", id).
    Scan(&user.ID, &user.Nombre, &user.Apellido, &user.Correo, &user.Password, &user.Tipo)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
    }

    return c.JSON(user)
}


// Actualizar usuario por ID
func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user := new(models.User)

    if err := c.BodyParser(user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    _, err := config.DB.Exec(context.Background(),
        "UPDATE usuarios SET nombre=$1, apellido=$2, correo=$3,  password=$4, tipo_usuario=$5 WHERE id_usuario=$6",
        user.Nombre, user.Apellido, user.Correo, user.Password, user.Tipo, id)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar usuario"})
    }

    return c.JSON(fiber.Map{"message": "Usuario actualizado correctamente"})
}


// Eliminar usuario por ID
func DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")

    _, err := config.DB.Exec(context.Background(),
        "DELETE FROM usuarios WHERE id_usuario=$1", id)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar usuario"})
    }

    return c.JSON(fiber.Map{"message": "Usuario eliminado correctamente"})
}


// Login permite a un usuario autenticarse y devuelve un token JWT
func Login(c *fiber.Ctx) error {
    var datos struct {
        Correo   string `json:"correo"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&datos); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // Buscar al usuario por correo
    var id int
    var password string
    err := config.DB.QueryRow(context.Background(),
    "SELECT id_usuario, password FROM usuarios WHERE correo=$1", datos.Correo).Scan(&id, &password)


    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Correo o contraseña incorrectos"})
    }

    /* // En esta versión no usamos hash, así que comparamos directo (luego agregaremos hash)
    if datos.Password != password {
        return c.Status(401).JSON(fiber.Map{"error": "Correo o contraseña incorrectos"})
    } */

    // Generamos el token
    token, err := utils.GenerarToken(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al generar token"})
    }

    return c.JSON(fiber.Map{"token": token})
}