package main

import (
	"fmt"
	"regexp"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"net/http"
)

// isEmail verifica si un string es un email válido
func IsEmail(email string) bool {
	// Expresión regular para validar emails
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	// Compila la expresión regular
	regex, err := regexp.Compile(emailRegex)
	if err != nil {
		fmt.Println("Error al compilar la expresión regular:", err)
		return false
	}

	// Verifica si el email coincide con la expresión regular
	return regex.MatchString(email)
}

// Función para validar y obtener los claims del token JWT
func ValidateAndGetClaims(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your_secret_key"), nil
    })
    if err != nil {
        return nil, fmt.Errorf("Error al validar el token: %w", err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("Token inválido")
}

func ValidateHeaders(c *gin.Context) (*Claims, error) {
    // Obtener el token JWT de la cabecera de la solicitud
    tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó un token de autenticación"})
        return nil, fmt.Errorf("no token provided")
    }

    // Validar y decodificar el token JWT
    claims, err := ValidateAndGetClaims(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return nil, err
    }

    return claims, nil
}