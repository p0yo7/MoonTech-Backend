// controllers.go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
)


// CreateUser maneja la solicitud POST para crear un nuevo usuario
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Create(&user)
    c.JSON(http.StatusOK, user)
}


func CreateProject(c *gin.Context) {
    // Validar los headers y obtener los claims
    claims, err := ValidateHeaders(c)
    if err != nil {
        return // Ya se manejó el error dentro de ValidateHeaders
    }

    // Crear un nuevo proyecto
    var proj Project
    if err := c.ShouldBindJSON(&proj); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Asignar el ID de usuario del JWT al campo Owner del proyecto
    proj.owner = int(claims.UserID)

    // Guardar el proyecto en la base de datos
    result := DB.Create(&proj)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
        fmt.Println(result.Error)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Project created successfully"})
}

func CreateRequirement(c *gin.Context){
    // Validar headers y obtener claims
    claims, err := ValidateHeaders(c)
    if err != nil {
        return
    }

    var req Requirements
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Asignar ID sacado del JWT
    req.owner = int(claims.UserID)
    
    // Crear en DB si es successfull
    result := DB.Create(&req)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create requirement"})
        fmt.Println(result.Error)
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Requirement created successfully"})
}