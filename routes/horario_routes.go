package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
    "hospital-back/middleware"
)

// HorarioRoutes que aagrupa los endpointss de loshorarios
func HorarioRoutes(app *fiber.App) {
    app.Post("/horarios",
        middleware.Autenticacion([]string{"crear_horarios"}),
        handlers.CreateHorario,
    )
    app.Get("/horarios",
        middleware.Autenticacion([]string{"ver_horarios"}),
        handlers.GetHorarios,
    )
    app.Get("/horarios/:id",
        middleware.Autenticacion([]string{"ver_horarios"}),
        handlers.GetHorarioByID,
    )
    app.Put("/horarios/:id",
        middleware.Autenticacion([]string{"actualizar_horarios"}),
        handlers.UpdateHorario,
    )
    app.Delete("/horarios/:id",
        middleware.Autenticacion([]string{"eliminar_horarios"}),
        handlers.DeleteHorario,
    )
}