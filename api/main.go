// main.go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "time"
)

func main() {
    // Conectar a la base de datos
    ConnectDatabase()

    // Mensaje de conexión exitosa a la base de datos
    log.Println("Connected to the database successfully.")

    // Crear una nueva instancia del router Gin
    r := gin.Default()

    // Configurar el middleware de CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Permitir todos los orígenes
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour, // Caché para las respuestas de preflight
    }))

    // Configurar las rutas
    SetupRouter(r)

    // Iniciar el servidor en el puerto 8080
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run the server: %v", err)
    }

    // Mensaje de inicio exitoso del servidor
    log.Println("Server started successfully on port 8080.")
}
