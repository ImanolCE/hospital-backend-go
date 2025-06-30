package routes

import (
    "hospital-back/handlers"
    "hospital-back/middleware"
    "github.com/gofiber/fiber/v2"
)

// UserRoutes define las rutas del CRUD de usuarios y el login
func UserRoutes(app *fiber.App) {
    app.Post("/usuarios", handlers.CreateUser)  // Registro
    app.Post("/login", handlers.Login)          // Login

    // Middleware protege lo que viene después
    app.Use(middleware.Autenticacion)

    app.Get("/usuarios", handlers.GetUsers)
    app.Get("/usuarios/:id", handlers.GetUserByID)
    app.Put("/usuarios/:id", handlers.UpdateUser)
    app.Delete("/usuarios/:id", handlers.DeleteUser)
}
