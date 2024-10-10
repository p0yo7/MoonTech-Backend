// auth.go
package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// rutas necesarias para hacer manejo de usuario y sus datos
var (
	microsoftEndpoint = oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}

	googleOAuth2Config = &oauth2.Config{
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback/google",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	microsoftOAuth2Config = &oauth2.Config{
		ClientID:     "YOUR_MICROSOFT_CLIENT_ID",
		ClientSecret: "YOUR_MICROSOFT_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback/microsoft",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     microsoftEndpoint,
	}

	store = sessions.NewCookieStore([]byte("secret-key"))
)

func handleCallback(config *oauth2.Config, c *gin.Context) {
	// El código de handleCallback permanece sin cambios
	// ...
}

func loginWithGoogleHandler(c *gin.Context) {
	url := googleOAuth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func loginWithMicrosoftHandler(c *gin.Context) {
	url := microsoftOAuth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func googleCallbackHandler(c *gin.Context) {
	handleCallback(googleOAuth2Config, c)
}

func microsoftCallbackHandler(c *gin.Context) {
	handleCallback(microsoftOAuth2Config, c)
}

// Estructura para el token
type Claims struct {
	UserID        uint   `json:"user_id"`
	Role          string `json:"role"`
	Team          int    `json:"team"` // Cambia a int si necesitas un ID de equipo
	UserFirstName string `json:"user_first_name"`
	jwt.RegisteredClaims
}

// Función para generar un JWT
func generateJWT(userID uint, role string, teamID int, userFirstName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Define el tiempo de expiración del token

	claims := &Claims{
		UserID:        userID,
		Role:          role,
		Team:          teamID, // Cambia a teamID
		UserFirstName: userFirstName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("your_secret_key") // Cambia esto por una clave secreta segura
	return token.SignedString(secret)
}

func Native_login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user Users
	result := DB.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Asegúrate de que user tiene un campo TeamID que es de tipo int
	teamID := user.TeamID // Cambia a TeamID si corresponde al ID del equipo
	userFirstName := user.FirstName

	// Generar el token JWT con el equipo y el primer nombre del usuario
	token, err := generateJWT(uint(user.ID), user.Role, teamID, userFirstName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Remover la contraseña antes de enviar los datos del usuario
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func CreateUser(c *gin.Context) {

	var user Users
	// Intenta bindear el user dado en el context al modelo de User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Println(err)
		return
	}

	//Verifica que sea un email valido
	if !IsEmail(user.WorkEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email lol get good "})
		return
	}

	// Verificar si el usuario ya existe
	var existingUser Users
	if err := DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		// Si no se encuentra un error, significa que el usuario ya existe
		fmt.Println("duplicate user found")
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate user found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	fmt.Println(hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	// Crear el nuevo usuario
	result := DB.Create(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func RegisterAuthRoutes(r *gin.Engine) {
	r.GET("/login/google", loginWithGoogleHandler)
	r.GET("/login/microsoft", loginWithMicrosoftHandler)
	r.GET("/callback/google", googleCallbackHandler)
	r.GET("/callback/microsoft", microsoftCallbackHandler)
	r.POST("/login/native", Native_login)
	r.POST("/createUser", CreateUser)
}
