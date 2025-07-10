
package main

import (
    "hospital-back/config" 
    "hospital-back/routes" 
    "hospital-back/middleware"

    "github.com/gofiber/fiber/v2" 
    "github.com/joho/godotenv"   

     "github.com/gofiber/fiber/v2/middleware/limiter"
     "time"
)

func main() {
	// Carga el archivo .env
    godotenv.Load() 

	// Conecta a Supabase
    config.Connect() 

	// Crea la instancia del servidor Fiber
    app := fiber.New() 

      // para limitar a 100 peticiones por minuto por su IP
    app.Use(limiter.New(limiter.Config{
        Max:        100,
        Expiration: 1 * time.Minute,
    }))

    // Para los logs 
    app.Use(middleware.Logger())


	// Registra las rutas de usuario
    routes.UserRoutes(app) 

    // Registra las rutas de los consultorios
     routes.ConsultorioRoutes(app)

    // Registra las rutas de los horarios
     routes.HorarioRoutes(app)

    // Regita las rutas de las consultas 
     routes.ConsultaRoutes(app)

	// Inicia el servidor en el puerto 3000
    app.Listen(":3000") 
}


