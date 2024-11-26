package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"bytes"
	"os"
	"io"
	"strings"
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
	proj.Status = 1
	StartDate := time.Now().Format("02-01-2006")
	parsedTime, err := time.Parse("02-01-2006", StartDate)
	if err != nil {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "Date parsing error"})
    	return
	}
	proj.StartDate = parsedTime

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
    ProjectID          int    `json:"project_id" gorm:"column:project_id"`
    ProjectName        string `json:"project_name" gorm:"column:project_name"`
    ProjectDescription string `json:"project_description" gorm:"column:project_description"`
    ProjectBudget      int    `json:"project_budget" gorm:"column:project_budget"`
    CompanyName        string `json:"company_name" gorm:"column:company_name"`
    CompanyDescription string `json:"company_description" gorm:"column:company_description"`
    CompanySize        int    `json:"company_size" gorm:"column:company_size"`
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
	RequirementID            int       `json:"requirement_id"`            // ID del proyecto al que pertenece
	RequirementDescription   string    `json:"requirement_description"`   // Descripción del requerimiento
	RequirementApproved      bool      `json:"requirement_approved"`      // Estado de aprobación del requerimiento
	RequirementTimeStamp	 time.Time `json:"requirement_timestamp"`     // Fecha y hora de creación del requerimiento

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
    TaskEstimatedCost int    `json:"task_estimated_cost"` // Costo estimado para la tarea
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
	ProjectID 		int    `json:"id"`
	ProjectName        string `json:"projName"`
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

// GetUsersWithTeamByName obtiene los usuarios que pertenecen a un equipo especificado por el nombre
func GetUsersWithTeamByName(c *gin.Context) {
    // Validar los encabezados y obtener los claims
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

    // Obtener el nombre del equipo desde los parámetros de la URL
    teamName := c.Param("team_name")

    // Crear un slice para almacenar los resultados
    var usersWithTeams []struct {
        ID       int   `json:"id"`
        Username string `json:"username"`
        TeamName string `json:"team_name"`
    }

    // Ejecutar la consulta SQL para obtener los usuarios y sus equipos filtrando por el nombre del equipo
    result := DB.Raw(`
        SELECT u.id, u.username, t.team_name
        FROM users u
        INNER JOIN teams t ON u.team = t.id
        WHERE t.team_name = ?
    `, teamName).Scan(&usersWithTeams)

    // Verificar si hubo errores al ejecutar la consulta
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users with teams"})
        return
    }

    // Retornar la respuesta en formato JSON
    c.JSON(http.StatusOK, usersWithTeams)
}

func InsertUserToProject(c *gin.Context) {
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

    // Crear una instancia del modelo ProjectUser
    var projectUser ProjectUser

    // Validar y enlazar datos JSON al modelo
    if err := c.ShouldBindJSON(&projectUser); err != nil {
        fmt.Println("Error al analizar JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    // Validar si el usuario y el proyecto existen antes de intentar insertarlos (opcional)
    var userExists, projectExists bool

    if err := DB.Model(&Users{}).
        Select("count(*) > 0").
        Where("id = ?", projectUser.UserID).
        Find(&userExists).Error; err != nil || !userExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        fmt.Println("Error validando usuario:", err)
        return
    }

    if err := DB.Model(&Projects{}).
        Select("count(*) > 0").
        Where("id = ?", projectUser.ProjectID).
        Find(&projectExists).Error; err != nil || !projectExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        fmt.Println("Error validando proyecto:", err)
        return
    }

    // Intentar crear el registro en la base de datos
    if err := DB.Create(&projectUser).Error; err != nil {
        fmt.Println("Error al insertar usuario en el proyecto:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user to project"})
        return
    }

    // Responder con éxito
    c.JSON(http.StatusOK, gin.H{
        "message": "User successfully added to project",
        "project_user": projectUser,
    })
}


// RequestBody estructura para el body de la petición
type GenerateRequirementsRequest struct {
    ProjectID int `json:"project_id" binding:"required"`
}

// OpenAIRequest estructura para la petición a OpenAI
type OpenAIRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Temperature float64   `json:"temperature"`
}

// Message estructura para los mensajes de OpenAI
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// OpenAIResponse estructura para la respuesta de OpenAI
type OpenAIResponse struct {
    ID      string `json:"id"`
    Object  string `json:"object"`
    Created int    `json:"created"`
    Model   string `json:"model"`
    Choices []struct {
        Index   int `json:"index"`
        Message struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
        FinishReason string `json:"finish_reason"`
    } `json:"choices"`
}


func GenerateAndCreateProjectRequirements(c *gin.Context) {
    // Validar headers y obtener claims
    claims, err := ValidateHeaders(c)
    if err != nil {
        if errors.Is(err, errors.New("Token expirado")) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        }
        return
    }
    // Obtener project_id del body
    var reqBody GenerateRequirementsRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Obtener información del proyecto usando el stored procedure
    var projectInfo ProjectInfo
    result := DB.Raw("CALL GetProjectInfo(?)", reqBody.ProjectID).Scan(&projectInfo)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project information"})
        fmt.Println(result.Error)
        return
    }

    // Verificar que el proyecto existe
    if projectInfo.ProjectID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }

    // Preparar el prompt para OpenAI
    prompt := fmt.Sprintf(
        "Basándote en la siguiente información del proyecto, genera una lista de requerimientos técnicos y funcionales.\n"+
            "Genera cada requerimiento en una línea separada, sin numeración ni viñetas.\n"+
            "Cada requerimiento debe ser una descripción clara y concisa.\n\n"+
            "Información del proyecto:\n"+
            "Nombre del Proyecto: %s\n"+
            "Descripción: %s\n"+
            "Presupuesto: %f\n"+
            "Empresa: %s\n"+
            "Descripción de la empresa: %s\n"+
            "Tamaño de la empresa: %d empleados\n\n"+
            "Por favor, genera los requerimientos considerando el tamaño de la empresa y el presupuesto disponible.",
        projectInfo.ProjectName,
        projectInfo.ProjectDescription,
        projectInfo.ProjectBudget,
        projectInfo.CompanyName,
        projectInfo.CompanyDescription,
        projectInfo.CompanySize,
    )

    // Preparar la petición a OpenAI
    openAIReq := OpenAIRequest{
        Model: "gpt-3.5-turbo",
        Messages: []Message{
            {
                Role:    "system",
                Content: "Eres un experto analista de sistemas que ayuda a definir requerimientos de proyectos de software. Genera requerimientos concisos y bien estructurados, uno por línea.",
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        Temperature: 0.7,
    }

    jsonData, err := json.Marshal(openAIReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare OpenAI request"})
        return
    }

    // Realizar la petición a OpenAI
    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
        return
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Failed to connect to OpenAI: %v", err),
        })
        return
    }
    defer resp.Body.Close()

    // Verificar el código de estado de la respuesta
    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("OpenAI API returned status %d", resp.StatusCode),
            "details": string(bodyBytes),
        })
        return
    }

    // Procesar la respuesta usando la estructura definida
    var openAIResponse OpenAIResponse
    if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Failed to parse OpenAI response: %v", err),
        })
        return
    }

    // Verificar que hay choices disponibles
    if len(openAIResponse.Choices) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No response content from OpenAI"})
        return
    }

    // Obtener el contenido de la respuesta
    content := openAIResponse.Choices[0].Message.Content

    // Verificar que el contenido no está vacío
    if content == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty response content from OpenAI"})
        return
    }

    // Procesar los requerimientos
    requirementsList := strings.Split(content, "\n")
    var createdReqs []Requirements

    for _, reqText := range requirementsList {
        reqText = strings.TrimSpace(reqText)
        if reqText == "" {
            continue
        }

        // Crear el requerimiento
        req := Requirements{
            ProjectID:              reqBody.ProjectID,
            OwnerID:               int(claims.UserID),
            RequirementDescription: reqText,
            timestamp:             time.Now(),
            approved:              false,
			ApproverID:            int(claims.UserID),
        }

        // Insertar en la base de datos
        result := DB.Create(&req)
        if result.Error != nil {
            fmt.Printf("Error creating requirement: %v\n", result.Error)
            continue
        }

        createdReqs = append(createdReqs, req)
    }

    // Retornar respuesta
    if len(createdReqs) > 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully created %d requirements", len(createdReqs)),
            "requirements": createdReqs,
        })
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create any requirements",
        })
    }
}
// RequestBody estructura para el body de la petición
type GenerateTasksRequest struct {
    ProjectID int `json:"project_id" binding:"required"`
}

func GenerateAndCreateTasks(c *gin.Context) {
    // Validar headers y obtener claims
    claims, err := ValidateHeaders(c)
    if err != nil {
        if errors.Is(err, errors.New("Token expirado")) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expirado"})
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        }
        return
    }

    // Obtener project_id del body
    var reqBody GenerateTasksRequest
    if err := c.ShouldBindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Obtener todos los requerimientos del proyecto
    var requirements []RequirementResponse
    result := DB.Raw("CALL GetProjectRequirements(?)", reqBody.ProjectID).Scan(&requirements)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching project requirements"})
        return
    }

    if len(requirements) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No requirements found for this project"})
        return
    }

    var allCreatedTasks []Tasks
    var failedRequirements []int

    // Iterar sobre cada requerimiento
    for _, requirement := range requirements {
        // Preparar el prompt para OpenAI
        prompt := fmt.Sprintf(
            "Basándote en el siguiente requerimiento, genera una lista de tareas técnicas necesarias.\n"+
                "Para cada tarea, especifica:\n"+
                "- Nombre de la tarea\n"+
                "- Descripción detallada\n"+
                "- Tiempo estimado en horas\n"+
                "- Costo estimado (basado en $50/hora)\n\n"+
                "Requerimiento: %s\n\n"+
                `Responde en formato JSON como este ejemplo:
                [
                    {
                        "name": "Nombre de la tarea",
                        "description": "Descripción detallada",
                        "estimatedTime": 8,
                        "estimatedCost": 400
                    }
                ]`,
            requirement.RequirementDescription,
        )

        // Preparar la petición a OpenAI
        openAIReq := OpenAIRequest{
            Model: "gpt-3.5-turbo",
            Messages: []Message{
                {
                    Role:    "system",
                    Content: "Eres un experto técnico que ayuda a desglosar requerimientos en tareas técnicas específicas y estimables.",
                },
                {
                    Role:    "user",
                    Content: prompt,
                },
            },
            Temperature: 0.7,
        }

        jsonData, err := json.Marshal(openAIReq)
        if err != nil {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            continue
        }

        // Realizar la petición a OpenAI
        req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
        if err != nil {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            continue
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            resp.Body.Close()
            continue
        }

        // Procesar la respuesta de OpenAI
        var openAIResponse OpenAIResponse
        if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            resp.Body.Close()
            continue
        }
        resp.Body.Close()

        // Verificar que hay contenido en la respuesta
        if len(openAIResponse.Choices) == 0 || openAIResponse.Choices[0].Message.Content == "" {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            continue
        }

        // Parsear el JSON de las tareas generadas
        type GeneratedTask struct {
            Name          string  `json:"name"`
            Description   string  `json:"description"`
            EstimatedTime int     `json:"estimatedTime"`
            EstimatedCost float64 `json:"estimatedCost"`
        }

        var generatedTasks []GeneratedTask
        if err := json.Unmarshal([]byte(openAIResponse.Choices[0].Message.Content), &generatedTasks); err != nil {
            failedRequirements = append(failedRequirements, requirement.RequirementID)
            continue
        }

        // Crear las tareas en la base de datos
        for _, genTask := range generatedTasks {
            task := Tasks{
                RequirementID:  requirement.RequirementID,
                Team:          1, // Valor por defecto
                CreatedBy:     int(claims.UserID),
                Name:          genTask.Name,
                Description:   genTask.Description,
                Language:      1, // Valor por defecto
                EstimatedTime: genTask.EstimatedTime,
                EstimatedCost: int(genTask.EstimatedCost),
                Ajuste:        0,
            }

            result := DB.Create(&task)
            if result.Error != nil {
                fmt.Printf("Error creating task for requirement %d: %v\n", requirement.RequirementID, result.Error)
                continue
            }

            allCreatedTasks = append(allCreatedTasks, task)
        }
    }

    // Preparar la respuesta
    response := gin.H{
        "message":     fmt.Sprintf("Successfully created %d tasks", len(allCreatedTasks)),
        "tasks":       allCreatedTasks,
        "total_tasks": len(allCreatedTasks),
    }

    if len(failedRequirements) > 0 {
        response["failed_requirements"] = failedRequirements
        response["warning"] = "Some requirements failed to generate tasks"
    }

    // Retornar respuesta
    if len(allCreatedTasks) > 0 {
        c.JSON(http.StatusOK, response)
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":               "Failed to create any tasks",
            "failed_requirements": failedRequirements,
        })
    }
}