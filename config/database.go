package config


import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
    
    url := os.Getenv("SUPABASE_URL")
    if url == "" {
        log.Fatal("SUPABASE_URL no está configurada")
    }

    // Si ya existe un pool abierto, lo cerramos antes de abrir otro
    if DB != nil {
        DB.Close()
    }

    var err error
    DB, err = pgxpool.New(context.Background(), url)
    if err != nil {
        log.Fatal("Error al conectar a Supabase:", err)
    }

    log.Println("Conexión a Supabase exitosa")
}
