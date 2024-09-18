// db.go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

// ConnectDatabase establece la conexión a la base de datos SQLite
func ConnectDatabase() {
    var err error
    DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    // Migrar el esquema
    DB.AutoMigrate(&User{})
}
