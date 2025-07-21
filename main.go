
package main

import (
    "hospital-back/config" 
    "hospital-back/routes" 
    "hospital-back/middleware"

    "github.com/gofiber/fiber/v2" 
    "github.com/joho/godotenv"   

     "github.com/gofiber/fiber/v2/middleware/limiter"
     "time"

     "github.com/gofiber/fiber/v2/middleware/cors"

      "log"
)

func main() {
    
	// 1) Carga variables de entorno desde .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, relying on environment variables")
    } 

	// Conecta a Supabase
    config.Connect() 

    // para cerrar la conexion al temrinar
    defer config.DB.Close()

	// Crea la instancia del servidor Fiber
    app := fiber.New() 

    app.Use(middleware.Logger()) 

    app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:4200", // o "*" para cualquier origen
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))


      // para limitar a 100 peticiones por minuto por su IP
    app.Use(limiter.New(limiter.Config{
        Max:        100,
        Expiration: 1 * time.Minute,
    }))

	// Registra las rutas de usuario
    routes.UserRoutes(app) 

    // Registra las rutas de los consultorios
     routes.ConsultorioRoutes(app)

    // Registra las rutas de los horarios
     routes.HorarioRoutes(app)

    // Regita las rutas de las consultas 
     routes.ConsultaRoutes(app)


	// Inicia el servidor en el puerto 3000
    log.Fatal(app.Listen(":3000"))

    

}


