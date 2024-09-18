// models.go
package main

import (
    "gorm.io/gorm"
)

// User representa el modelo de usuario
type User struct {
    gorm.Model
    Name  string `json:"name"`
    Email string `json:"email"`
}
