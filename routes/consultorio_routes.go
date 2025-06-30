package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
)

// ConsultorioRoutes agrupa los endpoints de los consultorios
func ConsultorioRoutes(app *fiber.App) {
    app.Post("/consultorios", handlers.CreateConsultorio)
    app.Get("/consultorios", handlers.GetConsultorios)
    app.Get("/consultorios/:id", handlers.GetConsultorioByID)
    app.Put("/consultorios/:id", handlers.UpdateConsultorio)
    app.Delete("/consultorios/:id", handlers.DeleteConsultorio)
}
