
package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
)


// UserRoutes define las rutas relacionadas con los usuarios
func UserRoutes(app *fiber.App) {
    app.Post("/usuarios", handlers.CreateUser) // para registrar un nuevo usuario
}
