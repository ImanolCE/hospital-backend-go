
package config

import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v5"
)


var DB *pgx.Conn

// Connect establece la conexión a Supabase usandoo la URL del archivo .env
func Connect() {
    url := os.Getenv("SUPABASE_URL") 
    if url == "" {
        log.Fatal("SUPABASE_URL no está configurada") 
    }

    var err error
    // aqui coonecta a la base de datos de Supabase
    DB, err = pgx.Connect(context.Background(), url)
    if err != nil {
        log.Fatal("Error al conectar a Supabase:", err)
    }

    log.Println("Conexión a Supabase exitosa") 
}

