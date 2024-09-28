// routers.go
package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "os"
)

// Configura el logger para que escriba en un archivo
func setupLogOutput() {
    file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("No se pudo abrir el archivo de log: %v", err)
    }
    gin.DefaultWriter = file
}

// SetupRouter configura las rutas de la API
func SetupRouter(r *gin.Engine) *gin.Engine {
    setupLogOutput()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.POST("/login/native", Native_login)
    r.POST("/createUser", CreateUser)
    r.POST("/createProject", CreateProject)
    r.POST("/createRequirement", CreateRequirement)
    // Dockers
    // Docker compose
    // Aprobar requerimiento
    // Rechazar requerimiento
    // Modificar requerimiento
    // Algoritmo de parentezco para contratos marco
    // Generar reporte
    // Recibir los proyectos de un usuario
    // Abrir un proyecto
    // Hacer llamada a llama3 para generar tareas
    // Enviar Feedback de tareas
    // Llamada para recibir tareas
    // Dashboard de vista summary de proyecto 
    // Ver lo de keys de microsoft y google auth
    // Notificaciones
    return r
}
