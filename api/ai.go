package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Estructura para los datos que recibe el webhook
type WebhookRequest struct {
	Message string `json:"message"`
}

// Estructura para los datos que se enviarán al servidor AI
type AIRequest struct {
	Input string `json:"input"`
}

// Estructura para la respuesta del servidor AI
type AIResponse struct {
	Result string `json:"result"`
}

// Función que maneja el webhook
func SendRequirementsAI(c *gin.Context) {
	var webhookReq WebhookRequest

	// Parsear el cuerpo de la solicitud en JSON
	if err := c.ShouldBindJSON(&webhookReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al parsear el cuerpo de la solicitud"})
		return
	}

	// Crear la solicitud para el servidor AI
	aiReq := AIRequest{
		Input: webhookReq.Message,
	}
	aiReqBody, err := json.Marshal(aiReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar la solicitud para el servidor AI"})
		return
	}

	// Realizar la solicitud POST al servidor AI
	// De momento el puerto es variable
	aiServerURL := "http://localhost:8501/getRequirements"
	resp, err := http.Post(aiServerURL, "application/json", bytes.NewBuffer(aiReqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al hacer la solicitud al servidor AI"})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta del servidor AI
	aiRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer la respuesta del servidor AI"})
		return
	}

	// Parsear la respuesta del servidor AI
	var aiResp AIResponse
	err = json.Unmarshal(aiRespBody, &aiResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear la respuesta del servidor AI"})
		return
	}

	// Devolver la respuesta al cliente que llamó al webhook
	c.JSON(http.StatusOK, aiResp)
}