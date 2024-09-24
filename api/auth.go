package main

import (
    // "encoding/json"
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "fmt"
)

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
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// Función para generar un JWT
func generateJWT(userID uint, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Define el tiempo de expiración del token

    claims := &Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := []byte("your_secret_key") // Cambia esto por una clave secreta segura
    return token.SignedString(secret)
}

func native_login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"` // Cambiado a Username y exportado
        Password string `json:"password"` // Cambiado para ser exportado
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var user User
    result := DB.Where("username = ?", loginData.Username).First(&user) // Cambiado a loginData.Username
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"}) // Mensaje actualizado
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        fmt.Println(err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"}) // Mensaje actualizado
        return
    }

    // Generar el token JWT
    token, err := generateJWT(user.ID, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    // Remover la contraseña antes de enviar los datos del usuario
    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Verificar si el usuario ya existe
    var existingUser User
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
    r.POST("/login/native", native_login)
    r.POST("/register", register)
}