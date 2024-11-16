package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetSchema(c *gin.Context) {
	// Obtener el nombre de la base de datos
	databaseName := DB.Migrator().CurrentDatabase()

	// Obtener todas las tablas usando una consulta SQL
	var tables []string
	rows, err := DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las tablas"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	// Responder con el nombre de la base de datos y las tablas
	c.JSON(http.StatusOK, gin.H{
		"database": databaseName,
		"tables":   tables,
	})
}

func GetProjects(c *gin.Context) {
	// Validar headers
	claims, err := ValidateHeaders(c)
	if err != nil {
		return // Ya se manejó el error dentro de ValidateHeaders
	}
	fmt.Println(claims)
	// Obtener id
	// Obtener proyectos abiertos del usuario
}

func GetProjectsByID(c *gin.Context) {
	// Validar Headers
	// Obtener Id del proyecto
	fmt.Println(c)
}

func CreateProject(c *gin.Context) {
	// Validar los headers y obtener los claims
	claims, err := ValidateHeaders(c)
	if err != nil {
		return // Ya se manejó el error dentro de ValidateHeaders
	}

	// Crear un nuevo proyecto
	var proj Projects
	if err := c.ShouldBindJSON(&proj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar el ID de usuario del JWT al campo Owner del proyecto
	proj.Owner.ID = int(claims.UserID)

	// Guardar el proyecto en la base de datos
	result := DB.Create(&proj)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		fmt.Println(result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project created successfully", "id": proj.ID})
}

func CreateRequirement(c *gin.Context) {
	// Validar headers y obtener claims
	claims, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}

	var req Requirements
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar ID sacado del JWT
	req.Owner.ID = int(claims.UserID)

	// Valores defaults
	req.approved = false
	// Si el requerimiento no tiene un aprobador, asignar el mismo
	req.Approver.ID = int(claims.UserID)

	// Crear en DB si es successfull
	result := DB.Create(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create requirement"})
		fmt.Println(result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Requirement created successfully"})
}

func ApproveRequirement(c *gin.Context) {
	// Validar headers y obtener claims
	claims, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return
	}

	// Obtener el ID del requerimiento de los parámetros de la URL
	reqID := c.Param("id")

	var req Requirements
	// Buscar el requerimiento por ID
	if err := DB.First(&req, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Requirement not found"})
		return
	}

	// Actualizar el campo de aprobado
	req.approved = true
	req.ApproverID = int(claims.UserID) // Asignar el usuario que aprobó el requerimiento

	// Guardar los cambios en la base de datos
	if err := DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve requirement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Requirement approved successfully"})
}

func CreateTeam(c *gin.Context) {
	var team Teams

	// Bind the incoming JSON to the team struct
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save the team to the database
	if err := DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create team"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Team created successfully", "team": team})
}

func CreateArea(c *gin.Context) {
	var area Areas
	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&area)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create area"})
		fmt.Println(result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Area created successfully"})
}

func CreateCompany(c *gin.Context) {
	var company Companies
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&company)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		fmt.Println(result.Error)
		return
	}
	fmt.Println(company)
	c.JSON(http.StatusOK, gin.H{"message": "Company created successfully", "company": company})
}

func CreateBusinessType(c *gin.Context) {
	var business BusinessTypes
	if err := c.ShouldBindJSON(&business); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&business)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create business type"})
		fmt.Println(result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Business type created successfully"})
}

func CreateRepresentative(c *gin.Context) {
	var rep Representatives
	if err := c.ShouldBindJSON(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&rep)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create representative"})
		fmt.Println(result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Representative created successfully"})
}

func AssignProjectLeaders(users []Users, projId int, c *gin.Context) {
	// claims, err := ValidateHeaders(c)
	// if err != nil {
	// 	// Verificar si el error es de token expirado
	// 	if errors.Is(err, errors.New("Token expirado")) {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
	// 	} else {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	}
	// 	return
	// }
	// for _, user := range users {
	// 	leader := Leaders {
	// 		UserID: int(user.ID),
	// 		ProjID: int(projId),
	// 	}
	// 	result := DB.Create(&leader)
	// 	if result.Error != nil {
	// 		return result.Error
	// 	}
	// }
	fmt.Println("TODO")
	// return nil
}

// Obtener users por team
// Obtener areas
func NotifyProjectCreation(users []Users) {
	fmt.Println("XD")
}

func NotifyProjectTurn(users []Users) {
	fmt.Println("XD")
}

type ProjectInfo struct {
	CompanyName        string `json:"company_name"`
	CompanyDescription string `json:"company_description"`
	CompanySize        int    `json:"company_size"`
	ProjectName        string `json:"project_name"`
	ProjectDescription string `json:"project_description"`
	ProjectBudget      int    `json:"project_budget"`
}

func GetProjectGeneralInfo(c *gin.Context) {
	// Obtener el ID del proyecto y convertirlo a entero
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var projectInfo ProjectInfo

	// Llamar al procedimiento almacenado con el ID convertido
	result := DB.Raw("CALL GetProjectInfo(?)", projectID).Scan(&projectInfo)

	// Verificar si hubo errores en la consulta
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project information"})
		fmt.Println(result.Error)
		return
	}

	// Retornar la información del proyecto en formato JSON
	c.JSON(http.StatusOK, projectInfo)
}

type RequirementResponse struct {
	ID                   int       `json:"id"`                    // ID del requerimiento
	ProjectID            int       `json:"project_id"`            // ID del proyecto al que pertenece
	RequirementText      string    `json:"requirement_text"`      // Descripción del requerimiento
	RequirementApproved  bool      `json:"requirement_approved"`  // Fecha de creación del requerimiento
	RequirementTimestamp time.Time `json:"requirement_timestamp"` // Fecha de actualización del requerimiento
}

func GetProjectRequirements(c *gin.Context) {
	// Validar headers y obtener claims
	_, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return
	}
	// Obtener el ID del proyecto de los parámetros de la URL
	ProjID := c.Param("id")

	// Crear un slice para almacenar los requerimientos
	var requirements []RequirementResponse

	// Llamar al procedimiento almacenado usando db.Raw()
	result := DB.Raw("CALL GetProjectRequirements(?)", ProjID).Scan(&requirements)

	// Verificar si hubo errores al ejecutar el query
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project requirements"})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, requirements)
}

type TaskResponse struct {
	TaskID            int    `json:"task_id"`             // ID de la tarea
	TaskTitle         string `json:"task_title"`          // Nombre de la tarea
	TaskDescription   string `json:"task_description"`    // Descripción de la tarea
	TaskEstimatedTime int    `json:"task_estimated_time"` // Tiempo estimado para la tarea (en horas)
}

func GetProjectTasks(c *gin.Context) {
	// Validar headers y obtener claims
	_, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return
	}

	// Obtener el ID del proyecto de los parámetros de la URL
	ProjID := c.Param("id")

	// Crear un slice para almacenar las tareas
	var tasks []TaskResponse // Asumiendo que tienes un struct TaskResponse para las tareas

	// Llamar al procedimiento almacenado para obtener las tareas del proyecto
	result := DB.Raw("CALL GetProjectTasks(?)", ProjID).Scan(&tasks)

	// Verificar si hubo errores al ejecutar el query
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project tasks"})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, tasks)
}
func CreateTasks(c *gin.Context) {
	// Validar headers y obtener claims
	_, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return
	}

	var task Tasks

	// Validar y enlazar datos JSON
	if err := c.ShouldBindJSON(&task); err != nil {
		fmt.Println("Error al analizar JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	// Validar si el requirement existe antes de intentar crear la tarea
	// var requirementExists bool
	// if err := DB.Model(&Requirements{}).
	// 	Select("count(*) > 0").
	// 	Where("id = ?", task.RequirementID).
	// 	Find(&requirementExists).Error; err != nil || !requirementExists {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requirement ID"})
	// 	fmt.Println(err)
	// 	return
	// }

	// Intentar crear la tarea en la base de datos
	if err := DB.Create(&task).Error; err != nil {
		fmt.Println("Error al crear tarea:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		fmt.Println(task)
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

type ProjectDescription struct {
	ProjectID            int    `json:"project_id"`             // ID de la tarea
	ProjectDescription         string `json:"project_description"`          // Nombre de la tarea
	Budget   string `json:"project_budget"`    // Descripción de la tarea
}


func UpdateProjectDescription(c *gin.Context) {
	// Validar headers y obtener claims
	_, err := ValidateHeaders(c)
	if err != nil {
		// Verificar si el error es de token expirado
		if errors.Is(err, errors.New("Token expirado")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
		return
	}

	// Parsear la solicitud JSON
	var updatedProject ProjectDescription
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Buscar el proyecto existente por ID
	var project Projects
	if err := DB.First(&project, updatedProject.ProjectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Actualizar los campos del proyecto
	project.ProjectDescription = updatedProject.ProjectDescription
	budgetInt, err := strconv.Atoi(updatedProject.Budget)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "error": "Failed to update project" })
	}
	project.Budget = budgetInt

	// Guardar los cambios en la base de datos
	if err := DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Responder con los datos actualizados
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
// ProjectResponse estructura para mapear los resultados del procedimiento almacenado
type ProjectResponse struct {
	ProjectID          int    `json:"project_id"`
	ProjectName        string `json:"project_name"`
	ProjectDescription string `json:"project_description"`
	ProjectBudget      int    `json:"project_budget"`
	CompanyName        string `json:"company_name"`
	CompanyDescription string `json:"company_description"`
	CompanySize        string `json:"company_size"`
}

// GetActiveProjectsForUserHandler obtiene los proyectos activos de un usuario
func GetActiveProjectsForUser(c *gin.Context) {
	// Obtener el ID del usuario de los parámetros de la URL
	userID := c.Param("user_id")

	// Validar que el userID sea un entero válido
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Crear un slice para almacenar la respuesta del procedimiento almacenado
	var projects []ProjectResponse

	// Llamar al procedimiento almacenado con GORM
	result := DB.Raw("CALL GetActiveProjectsForUser(?)", userID).Scan(&projects)

	// Verificar si hubo errores en la ejecución
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects", "details": result.Error.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, projects)
}