package routes

import (
    "hospital-back/handlers"
    "hospital-back/middleware"
    "github.com/gofiber/fiber/v2"
)

 func UserRoutes(app *fiber.App) {
    app.Post("/usuarios", handlers.CreateUser)
    app.Post("/login", handlers.Login)
    app.Post("/refresh", handlers.Refresh)

    app.Post("/mfa/recovery", handlers.RecoverMfaStart)

    app.Post("/mfa/recovery/activate", handlers.ActivateMfaRecovery)


    app.Post("/mfa/activar", middleware.Autenticacion([]string{}), handlers.ActivarMFA)

    app.Get("/mfa/setup", middleware.Autenticacion([]string{}), handlers.MFASetup)
    app.Post("/mfa/verify", middleware.Autenticacion([]string{}), handlers.MFAVerify)

    app.Post("/mfa/regenerate",middleware.Autenticacion([]string{}),handlers.RegenerateMfa)

   
    app.Get("/usuarios",
        middleware.Autenticacion([]string{"ver_usuarios"}), handlers.GetUsers,
    )
    
    app.Get("/usuarios/:id",
        middleware.Autenticacion([]string{"ver_usuarios"}),
        handlers.GetUserByID,
    )
    
    app.Put("/usuarios/:id",
        middleware.Autenticacion([]string{"actualizar_usuario"}),
        handlers.UpdateUser,
    )
    app.Delete("/usuarios/:id",
        middleware.Autenticacion([]string{"eliminar_usuario"}),
        handlers.DeleteUser,
    )
}