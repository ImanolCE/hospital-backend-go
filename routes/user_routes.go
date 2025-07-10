package routes

import (
    "hospital-back/handlers"
    "hospital-back/middleware"
    "github.com/gofiber/fiber/v2"
)

// UserRoutes define las rutas del CRUD de usuarios y el login
/* func UserRoutes(app *fiber.App) {
    app.Post("/usuarios", handlers.CreateUser)  // Registro
    app.Post("/login", handlers.Login)          // Login

    // Middleware protege lo que viene después
    //app.Use(middleware.Autenticacion)

    // UserRoutes agrupa los endpoints de los usuarios
    app.Get("/usuarios", middleware.Autenticacion([]string{"ver_usuarios"}), handlers.GetUsers)
    app.Get("/usuarios/:id", middleware.Autenticacion([]string{"ver_usuarios"}), handlers.GetUserByID)
    app.Put("/usuarios/:id", middleware.Autenticacion([]string{"actualizar_usuario"}), handlers.UpdateUser)
    app.Delete("/usuarios/:id", middleware.Autenticacion([]string{"eliminar_usuario"}), handlers.DeleteUser)
}
 */

 func UserRoutes(app *fiber.App) {
    app.Post("/usuarios", handlers.CreateUser)
    app.Post("/login", handlers.Login)

     app.Post("/refresh", handlers.Refresh)

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