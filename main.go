
package main

import (
    "hospital-back/config" 
    "hospital-back/routes" 
    "github.com/gofiber/fiber/v2" 
    "github.com/joho/godotenv"   
)

func main() {
	// Carga el archivo .env
    godotenv.Load() 

	// Conecta a Supabase
    config.Connect() 

	// Crea la instancia del servidor Fiber
    app := fiber.New() 

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


