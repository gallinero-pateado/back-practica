package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRequest estructura de los datos recibidos
type practicasRequest struct {
	Titulo             string    `json:"Titulo"`
	Descripcion        string    `json:"Descripcion"`
	Id_empresa         int       `json:"Id_Empresa"`
	Ubicacion          string    `json:"Ubicacion"`
	Fecha_inicio       time.Time `json:"Fecha_inicio"`
	Fecha_fin          time.Time `json:"Fecha_fin"`
	Requisitos         string    `json:"Requisitos"`
	Fecha_expiracion   time.Time `json:"Fecha_expiracion"`
	Id_estado_practica int       `json:"Id_estado_practica"`
	Fecha_publicacion  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func UpdatePractica(c *gin.Context) {
	var req practicasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID de la práctica a actualizar desde los parámetros de la ruta
	id := c.Param("id")

	var practica models.Practica
	// Buscar la práctica en la base de datos
	if result := database.DB.First(&practica, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Práctica no encontrada"})
		return
	}
	localTime := time.Now().Local()

	// Actualizar los campos de la práctica con los datos de la solicitud
	practica.Titulo = req.Titulo
	practica.Descripcion = req.Descripcion
	practica.Id_empresa = req.Id_empresa
	practica.Ubicacion = req.Ubicacion
	practica.Fecha_inicio = req.Fecha_inicio
	practica.Fecha_fin = req.Fecha_fin
	practica.Requisitos = req.Requisitos
	practica.Fecha_expiracion = req.Fecha_expiracion
	practica.Fecha_publicacion = localTime

	// Guardar los cambios en la base de datos
	if result := database.DB.Save(&practica); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la práctica en la base de datos"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "La práctica fue actualizada exitosamente", "id_practica": practica.Id})
}
