package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
)

// HorarioRoutes que aagrupa los endpointss de loshorarios
func HorarioRoutes(app *fiber.App) {
    app.Post("/horarios", handlers.CreateHorario)
    app.Get("/horarios", handlers.GetHorarios)
    app.Get("/horarios/:id", handlers.GetHorarioByID)
    app.Put("/horarios/:id", handlers.UpdateHorario)
    app.Delete("/horarios/:id", handlers.DeleteHorario)
}
