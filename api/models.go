package main

import (
	"time"
)

// BusinessTypes representa la tabla de tipos de negocio
type BusinessTypes struct {
	BusinessTypeID int    `gorm:"column:business_type_id;primaryKey;autoIncrement"`
	BusinessName   string `gorm:"column:name;size:100"`
	BusinessColor  string `gorm:"column:color;size:20"`
}

func (BusinessTypes) TableName() string {
	return "businessTypes"
}

// Representatives representa la tabla de representantes
type Representatives struct {
	RepresentativeID          int    `gorm:"column:representative_id;primaryKey;autoIncrement"`
	representative_first_name string `gorm:"column:firstname;size:50"`
	representative_last_name  string `gorm:"column:lastname;size:50"`
	representative_work_email string `gorm:"column:work_email;size:100"`
	representative_work_phone string `gorm:"column:work_phone;size:20"`
}

func (Representatives) TableName() string {
	return "representatives"
}

// Areas representa la tabla de áreas
type Areas struct {
	AreaID           int    `gorm:"column:area_id;primaryKey;autoIncrement"`
	area_name        string `gorm:"column:area_name;size:100"`
	area_description string `gorm:"column:area_description;type:text"`
}

func (Areas) TableName() string {
	return "areas"
}

// Companies representa la tabla de compañías
type Companies struct {
	CompanyID        int             `gorm:"column:company_id;primaryKey;autoIncrement"`
	company_name     string          `gorm:"column:name;size:100"`
	RepresentativeID int             `gorm:"column:representativeId"`
	BusinessTypeID   int             `gorm:"column:businessType"`
	representative   Representatives `gorm:"foreignKey:RepresentativeID;references:RepresentativeID"`
	business_type    BusinessTypes   `gorm:"foreignKey:BusinessTypeID;references:BusinessTypeID"`
}

func (Companies) TableName() string {
	return "companies"
}

// Users representa la tabla de usuarios
type Users struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string `gorm:"column:username;size:50;not null"`
	FirstName string `gorm:"column:first_name;size:50"` // Asegúrate de que este campo esté correctamente definido
	LastName  string `gorm:"column:lastname;size:50"`
	WorkEmail string `gorm:"column:work_email;size:100;unique"`
	WorkPhone string `gorm:"column:work_phone;size:20"`
	Password  string `gorm:"column:password;size:255;not null"`
	AreaID    int    `gorm:"column:area"`     // Cambiado a AreaID
	LeaderID  *int   `gorm:"column:leaderId"` // Cambiado a LeaderID
	Position  string `gorm:"column:position;size:100"`
	Role      string `gorm:"column:role;size:50"`
	Area      Areas  `gorm:"foreignKey:AreaID;references:ID"` // Asegúrate que el foreignKey y references sean correctos
	Leader    *Users `gorm:"foreignKey:LeaderID;references:ID"`
}

// TableName establece el nombre de la tabla en la base de datos
func (Users) TableName() string {
	return "users" // Asegúrate de usar el nombre en minúsculas
}

// Projects representa la tabla de proyectos
type Projects struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	project_name string    `gorm:"column:proj_name;size:100"`
	OwnerID      int       `gorm:"column:owner"`
	CompanyID    int       `gorm:"column:company"`
	AreaID       int       `gorm:"column:area"`
	start_date   time.Time `gorm:"column:startDate;type:date"`
	Owner        Users     `gorm:"foreignKey:OwnerID;references:ID"`
	Company      Companies `gorm:"foreignKey:CompanyID;references:CompanyID"`
	Area         Areas     `gorm:"foreignKey:AreaID;references:AreaID"`
}

func (Projects) TableName() string {
	return "projects"
}

// Leaders representa la tabla de líderes
type Leaders struct {
	ID      int      `gorm:"column:id;primaryKey;autoIncrement"`
	ProjID  int      `gorm:"column:proj_id"`
	UserID  int      `gorm:"column:user_id"`
	AreaID  int      `gorm:"column:area_id"`
	Project Projects `gorm:"foreignKey:ProjID;references:ID"`
	User    Users    `gorm:"foreignKey:UserID;references:ID"`
	Area    Areas    `gorm:"foreignKey:AreaID;references:AreaID"`
}

func (Leaders) TableName() string {
	return "leaders"
}

// Requirements representa la tabla de requerimientos
type Requirements struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	ProjectID  int       `gorm:"column:projectId"`
	OwnerID    int       `gorm:"column:owner"`
	text       string    `gorm:"column:text;type:text"`
	timestamp  time.Time `gorm:"column:timestamp;type:datetime;default:CURRENT_TIMESTAMP"`
	approved   bool      `gorm:"column:approved"`
	ApproverID int       `gorm:"column:approverId"`
	Project    Projects  `gorm:"foreignKey:ProjectID;references:ID"`
	Owner      Users     `gorm:"foreignKey:OwnerID;references:ID"`
	Approver   Users     `gorm:"foreignKey:ApproverID;references:ID"`
}

func (Requirements) TableName() string {
	return "requirements"
}

// Comments representa la tabla de comentarios
type Comments struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	OwnerID   int       `gorm:"column:owner_id"`
	ParentID  *int      `gorm:"column:parent_id"` // Puede ser nulo
	Text      string    `gorm:"column:text;type:text"`
	Timestamp time.Time `gorm:"column:timestamp;type:datetime;default:CURRENT_TIMESTAMP"`
	Owner     Users     `gorm:"foreignKey:OwnerID;references:ID"`
}

func (Comments) TableName() string {
	return "comments"
}

// Tasks representa la tabla de tareas
type Tasks struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement"`
	AreaID        int       `gorm:"column:area_id"`
	Title         string    `gorm:"column:title;size:100"`
	CreatedBy     int       `gorm:"column:created_by"`
	Description   string    `gorm:"column:description;type:text"`
	Timestamp     time.Time `gorm:"column:timestamp;type:datetime;default:CURRENT_TIMESTAMP"`
	EstimatedTime int       `gorm:"column:estimated_time"`
	Approved      bool      `gorm:"column:approved"`
	ApproverID    int       `gorm:"column:approver_id"`
	Area          Areas     `gorm:"foreignKey:AreaID;references:AreaID"`
	Creator       Users     `gorm:"foreignKey:CreatedBy;references:ID"`
	Approver      Users     `gorm:"foreignKey:ApproverID;references:ID"`
}

func (Tasks) TableName() string {
	return "tasks"
}
