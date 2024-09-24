package main

import (
    "github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas de la API
func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.POST("/createProject", CreateProject)
    r.POST("/createRequirement", CreateRequirement)
    // Registrar las rutas de autenticación
    RegisterAuthRoutes(r)

    return r
}
