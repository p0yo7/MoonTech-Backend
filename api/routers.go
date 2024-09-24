package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "os"
)

// Configura el logger para que escriba en un archivo
func setupLogOutput() {
    // Abrir (o crear) un archivo de log
    file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("No se pudo abrir el archivo de log: %v", err)
    }

    // Redirigir el logger de Gin para escribir en el archivo
    gin.DefaultWriter = file

}

// SetupRouter configura las rutas de la API
func SetupRouter() *gin.Engine {
    // Configurar la salida del log antes de inicializar el router
    setupLogOutput()

    // Crear la instancia de Gin con el middleware Logger y Recovery
    r := gin.New()

    // Middleware para loguear cada request
    r.Use(gin.Logger())
    r.Use(gin.Recovery()) // Para manejar los errores y evitar que el servidor caiga

    // Rutas
    r.POST("/createProject", CreateProject)
    r.POST("/createRequirement", CreateRequirement)

    // Registrar las rutas de autenticación
    RegisterAuthRoutes(r)

    return r
}