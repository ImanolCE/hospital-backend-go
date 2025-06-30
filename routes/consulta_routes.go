package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
)

// ConsultaRoutes agrupa los endpoints de consultas
func ConsultaRoutes(app *fiber.App) {
    app.Post("/consultas", handlers.CreateConsulta)
    app.Get("/consultas", handlers.GetConsultas)
    app.Get("/consultas/:id", handlers.GetConsultaByID)
    app.Put("/consultas/:id", handlers.UpdateConsulta)
    app.Delete("/consultas/:id", handlers.DeleteConsulta)
}
