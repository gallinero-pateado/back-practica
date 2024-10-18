package Postulacionespractica

import (
	"fmt"
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type PostulacionRequest struct {
	Mensaje string `json:"mensaje"` // Campo para el mensaje del usuario
}

func Postularpractica(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la práctica de la URL
	practicaidStr := c.Param("practicaid")
	if practicaidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
		return
	}

	// Convertir practicaidStr a int
	var practicaid int
	if _, err := fmt.Sscan(practicaidStr, &practicaid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
		return
	}

	// Obtener el mensaje del cuerpo de la solicitud
	var req PostulacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar el usuario por el uid de Firebase
	var usuario models.Usuario
	result := database.DB.Where("firebase_usuario = ?", uid).First(&usuario)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario en la base de datos"})
		return
	}

	// Guardar el ID del usuario en una variable
	idUsuario := usuario.Id // idUsuario es de tipo uint

	// Verificar si ya existe una postulacion para el usuario y la práctica
	var existingPostulacion models.Postulacion
	err := database.DB.Where("id_usuario = ? AND id_practica = ?", idUsuario, practicaid).First(&existingPostulacion).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Ya has postulado a esta práctica"})
		return
	}

	localTime := time.Now().Local()

	// Si no existe, crear una nueva postulacion
	postulacion := models.Postulacion{
		Id_usuario:            int(idUsuario), // Conversión de uint a int
		Fecha_postulacion:     localTime,
		Id_estado_postulacion: 1,           // Ajusta según tu lógica de estado
		Mensaje:               req.Mensaje, // Mensaje proporcionado por el usuario
		Id_practica:           practicaid,
	}

	if err := database.DB.Create(&postulacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la postulación"})
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Postulación creada exitosamente"})
}
