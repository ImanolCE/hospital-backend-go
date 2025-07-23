package routes

import (
    "hospital-back/handlers"
    "github.com/gofiber/fiber/v2"
    "hospital-back/middleware"
)

// ConsultaRoutes agrupa los endpoints de consultas
func ConsultaRoutes(app *fiber.App) {
    app.Post("/consultas",
        middleware.Autenticacion([]string{"crear_consultas"}),
        handlers.CreateConsulta,
    )
     app.Get("/consultas/paciente/:id",
        middleware.Autenticacion([]string{"ver_consultas"}),
        handlers.GetConsultasByPaciente,  // <— tu nueva función
    )
    app.Get("/consultas",
        middleware.Autenticacion([]string{"ver_consultas"}),
        handlers.GetConsultas,
    )
    app.Get("/consultas/:id",
        middleware.Autenticacion([]string{"ver_consultas"}),
        handlers.GetConsultaByID,
    )
    app.Put("/consultas/:id",
        middleware.Autenticacion([]string{"actualizar_consultas"}),
        handlers.UpdateConsulta,
    )
    app.Delete("/consultas/:id",
        middleware.Autenticacion([]string{"eliminar_consultas"}),
        handlers.DeleteConsulta,
    )
}
