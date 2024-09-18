// models.go
package main

// import (
//     "gorm.io/gorm"
// )

// User representa la tabla de usuarios
type User struct {
    ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Username   string `gorm:"size:50;not null" json:"username"`
    Firstname  string `gorm:"size:50" json:"firstname"`
    Lastname   string `gorm:"size:50" json:"lastname"`
    WorkEmail  string `gorm:"size:100;unique" json:"workEmail"`
    WorkPhone  string `gorm:"size:20" json:"workPhone"`
    Password   string `gorm:"size:255;not null" json:"password"`
    Area       int    `json:"area"`
    LeaderID   int    `json:"leaderId"`
    Position   string `gorm:"size:100" json:"position"`
    Role       string `gorm:"size:50" json:"role"`
}

// Project representa la tabla de proyectos
type Project struct {
    ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    ProjName  string `gorm:"size:100" json:"projName"`
    Owner     int    `json:"owner"`  // Hace referencia al userId
    Company   int    `json:"company"` // Hace referencia a la tabla de Companies
    Area      int    `json:"area"`   // Hace referencia a la tabla de Areas
    StartDate string `json:"startDate"` // Tipo DATE
}

// Leader representa la tabla de líderes
type Leader struct {
    ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
    ProjID int  `json:"projId"` // Hace referencia a Project
    UserID int  `json:"userId"` // Hace referencia a User
    Area   int  `json:"area"`   // Hace referencia a Area
}

// Requirement representa la tabla de requerimientos
type Requirement struct {
    ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    ProjectID   int    `json:"projectId"` // Hace referencia a Project
    Owner       int    `json:"owner"`     // Hace referencia a User
    Text        string `json:"text"`
    Timestamp   string `json:"timestamp"` // Tipo DATETIME
    Approved    bool   `json:"approved"`
    ApproverID  int    `json:"approverId"` // Hace referencia a User
}

// Comment representa la tabla de comentarios
type Comment struct {
    ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Owner   int    `json:"owner"` // Hace referencia a User
    Parent  int    `json:"parent"` // Puede ser un requerimiento o una tarea
    Text    string `json:"text"`
    Timestamp string `json:"timestamp"` // Tipo DATETIME
}

// Task representa la tabla de tareas
type Task struct {
    ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Area         int    `json:"area"` // Hace referencia a Area
    Title        string `gorm:"size:100" json:"title"`
    CreatedBy    int    `json:"createdBy"` // Hace referencia a User
    Description  string `json:"description"`
    Timestamp    string `json:"timestamp"` // Tipo DATETIME
    EstimatedTime int    `json:"estimatedTime"`
    Approved      bool   `json:"approved"`
    ApproverID    int    `json:"approverId"` // Hace referencia a User
}

// Company representa la tabla de compañías
type Company struct {
    ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
    Name           string `gorm:"size:100" json:"name"`
    RepresentativeID int   `json:"representativeId"` // Hace referencia a la tabla de Representantes
    BusinessTypeID  int   `json:"businessType"` // Hace referencia a la tabla businessType
}

// BusinessType representa la tabla de tipos de negocio
type BusinessType struct {
    ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Name  string `gorm:"size:100" json:"name"`
    Color string `gorm:"size:20" json:"color"`
}

// Representative representa la tabla de representantes
type Representative struct {
    ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Firstname string `gorm:"size:50" json:"firstname"`
    Lastname  string `gorm:"size:50" json:"lastname"`
    WorkEmail string `gorm:"size:100" json:"workEmail"`
    WorkPhone string `gorm:"size:20" json:"workPhone"`
}

// Area representa la tabla de áreas
type Area struct {
    ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string `gorm:"size:100" json:"name"`
    Description string `json:"description"`
}
