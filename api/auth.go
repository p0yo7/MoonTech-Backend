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

func native_login(c *gin.Context) {
    var loginData struct {
        WorkEmail string `json:"workEmail"`
        Password  string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var user User
    result := DB.Where("work_email = ?", loginData.WorkEmail).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    session, _ := store.Get(c.Request, "session")
    session.Values["user_id"] = user.ID
    session.Save(c.Request, c.Writer)

    // Remove the password before sending the user data
    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = string(hashedPassword)

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