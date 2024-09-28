// db.go
package main

import (
    "log"
    "os"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
)

var DB *gorm.DB

// ConnectDatabase establece la conexión a la base de datos MySQL
func ConnectDatabase() {
    // Cargar variables de entorno desde el archivo .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Obtener valores de las variables de entorno
    user := os.Getenv("user")
    password := os.Getenv("password")
    database := os.Getenv("database")
    host := os.Getenv("host")
    port := os.Getenv("port")

    // Formar el Data Source Name (DSN)
    dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
    
    maxRetries := 10
    retryInterval := 5 * time.Second

    for i := 0; i < maxRetries; i++ {
        // Conectar a la base de datos MySQL
        DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Printf("Error al conectar a la base de datos: %v. Intento %d/%d", err, i+1, maxRetries)
            time.Sleep(retryInterval)
            continue
        }

        // Migrar el esquema
        err = DB.AutoMigrate(&User{})
        if err != nil {
            log.Fatalf("Failed to migrate the database schema: %v", err)
        }

        log.Println("Conexión exitosa a la base de datos")
        return
    }

    log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}
