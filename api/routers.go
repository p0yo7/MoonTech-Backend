// routers.go
package main

import (
	"log"
	"os"

	// "api/controllers"
	"github.com/gin-gonic/gin"
)

// Configura el logger para que escriba en un archivo
func setupLogOutput() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("No se pudo abrir el archivo de log: %v", err)
	}
	gin.DefaultWriter = file
}

// SetupRouter configura las rutas de la API
func SetupRouter(r *gin.Engine) *gin.Engine {
	setupLogOutput()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// r.POST

	// NATIVE LOGIN
	// Agregarle el nombre del usuario, rol, equipo y id
	r.POST("/login/native", Native_login)

	// CREATE USER
	r.POST("/createUser", CreateUser)

	// CREATE PROJECT
	r.POST("/createProject", CreateProject)

	// CREATE PROJECT_USERS
	r.POST("/createProjectUsers", InsertUserToProject)
	// CREATE REQUIREMENT
	r.POST("/createRequirement", CreateRequirement)

	// CREATE BUSINESS TYPE
	r.POST("/createBusinessType", CreateBusinessType)

	// CREATE REPRESENTATIVE
	r.POST("/createRepresentative", CreateRepresentative)

	// CREATE AREA
	r.POST("createArea", CreateArea)

	// CREATE COMPANY
	r.POST("/createCompany", CreateCompany)

	r.POST("/createTeam", CreateTeam)
	// ADD COMMENT
	// r.POST("/createComment", CreateComment)
	// Create project's tasks
	r.POST("/createTasks", CreateTasks)
	// SEND DATA TO META SERVER
	r.POST("/sendRequirements", SendRequirementsAI)
	r.POST("/Project/Requirements/Generate", GenerateAndCreateProjectRequirements)
	
	r.POST("/Project/Tasks/Generate", GenerateAndCreateTasks)
	
	// PUT

	// APPROVE REQUIREMENT
	r.PUT("/approveRequirement", ApproveRequirement)
	r.PUT("/ProjectDescription/:id", UpdateProjectDescription)
	// REJECT REQUIREMENT
	// r.PUT("/rejectRequirement", RejecRequirement)

	// MODIFY REQUIREMENT
	// r.PUT("/modifyRequirement", ModifyRequirement)

	// GETS
	r.GET("/getSchema", GetSchema)
	r.GET("/GetUsersForTeams/:team_name", GetUsersWithTeamByName)
	r.GET("/getAITasks")
	// comercial, lider digital, legal, equipo digital, pm, finanzas, gdm
	r.GET("/getTeams")  //Para la parte de teams
	// Get active projects for user 
	r.GET("/ActiveProjects/:user_id", GetActiveProjectsForUser)
	// Get a project's specific general information to use in AI
	r.GET("/Project/:id", GetProjectGeneralInfo)
	// Get a project's specific requirements as an array of objects
	r.GET("/Project/:id/Requirements", GetProjectRequirements)
	// Get project's specific requirements
	r.GET("/Project/:id/Tasks", GetProjectTasks)

	// Rechazar requerimiento
	// Modificar requerimiento
	// Algoritmo de parentezco para contratos marco
	// Generar reporte
	// Recibir los proyectos de un usuario
	// Abrir un proyecto
	// Hacer llamada a llama3 para generar tareas
	// Enviar Feedback de tareas
	// Llamada para recibir tareas
	// Dashboard de vista summary de proyecto
	// Ver lo de keys de microsoft y google auth
	// Notificaciones
	return r
}
