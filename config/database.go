package config


import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
    supabaseURL := os.Getenv("SUPABASE_URL")
    if supabaseURL == "" {
        log.Fatal("SUPABASE_URL no está definida en el entorno o .env")
    }

    cfg, err := pgxpool.ParseConfig(supabaseURL)
    if err != nil {
        log.Fatalf("Error parseando SUPABASE_URL: %v", err)
    }

    cfg.ConnConfig.PreferSimpleProtocol = true

    pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
    if err != nil {
        log.Fatalf("Error conectando a la base de datos: %v", err)
    }

    DB = pool
    log.Println("✅ Conexión exitosa a la base de datos Supabase")
}


/* func Connect() {
    // 1) Parsear configuración desde URL
    cfg, err := pgxpool.ParseConfig(os.Getenv("SUPABASE_URL"))
    if err != nil {
        log.Fatal("Error parseando SUPABASE_URL:", err)
    }

    // 2) Desactivar cache de prepared statements
    cfg.ConnConfig.PreferSimpleProtocol = true

    // 3) Conectar usando esa configuración
    pool, err := pgxpool.ConnectConfig(context.Background(), cfg)
    if err != nil {
        log.Fatal("Error conectando a la base de datos:", err)
    }

    // 4) Asignar al pool global
    DB = pool
}
 */



