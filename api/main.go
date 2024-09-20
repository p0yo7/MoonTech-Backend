// main.go
package main

import (
    "log"
)

func main() {
    // Conectar a la base de datos
    ConnectDatabase()

    // Mensaje de conexión exitosa a la base de datos
    log.Println("Connected to the database successfully.")

    // Configurar las rutas
    r := SetupRouter()

    // Iniciar el servidor en el puerto 5000
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run the server: %v", err)
    }

    // Mensaje de inicio exitoso del servidor
    log.Println("Server started successfully on port 8080.")
}
