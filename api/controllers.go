// controllers.go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
)

// GetUsers maneja la solicitud GET para obtener usuarios
func GetUsers(c *gin.Context) {
    var users []User
    if err := DB.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// CreateUser maneja la solicitud POST para crear un nuevo usuario
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    DB.Create(&user)
    c.JSON(http.StatusCreated, user)
}

func CreateProject(c *gin.Context){
    var proj Project
    if err := c.ShouldBindJSON(&proj); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := DB.Create(&proj)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
        fmt.Println(result.Error)
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func CreateRequirement(c *gin.Context){
    var req Requirements
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    result := DB.Create(&req)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create requirement"})
        fmt.Println(result.Error)
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Requirement created successfully"})
}