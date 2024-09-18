// controllers.go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
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
