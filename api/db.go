// db.go
package main

import (
    "log"
    "os"

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
        log.Fatalf("Error loading .env file")
    }

    // Obtener valores de las variables de entorno
    user := os.Getenv("user")
    password := os.Getenv("password")
    database := os.Getenv("database")
    host := os.Getenv("host")
    port := os.Getenv("port")

    // Formar el Data Source Name (DSN)
    dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
    
    // Conectar a la base de datos MySQL
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    // Migrar el esquema
    DB.AutoMigrate(&User{})
}
