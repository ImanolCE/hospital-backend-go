// routes/consultorio_routes.go
package routes


import (

    "github.com/gofiber/fiber/v2"
    "hospital-back/handlers"
    "hospital-back/middleware"

)

// ConsultorioRoutes agrupa los endpoints de los consultorios
func ConsultorioRoutes(app *fiber.App) {
    
    app.Post("/consultorios",
        middleware.Autenticacion([]string{"crear_consultorios"}),
        handlers.CreateConsultorio,
    )
    app.Get("/consultorios",
        middleware.Autenticacion([]string{"ver_consultorios"}),
        handlers.GetConsultorios,
    )
    app.Get("/consultorios/:id",
        middleware.Autenticacion([]string{"ver_consultorios"}),
        handlers.GetConsultorioByID,
    )
    app.Put("/consultorios/:id",
        middleware.Autenticacion([]string{"actualizar_consultorios"}),
        handlers.UpdateConsultorio,
    )
    app.Delete("/consultorios/:id",
        middleware.Autenticacion([]string{"eliminar_consultorios"}),
        handlers.DeleteConsultorio,
    )
}
