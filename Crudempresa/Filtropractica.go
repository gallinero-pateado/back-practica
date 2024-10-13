package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// obtiene las prácticas aplicando filtros opcionales
func FiltroPracticas(c *gin.Context) {
	// Obtener los parámetros de la solicitud (query parameters)
	modalidad := c.Query("modalidad")
	areaPractica := c.Query("area_practica")
	jornada := c.Query("jornada")
	ubicacion := c.Query("ubicacion")
	fecha := c.Query("fecha_publicacion")

	// Crear una lista de prácticas para almacenar los resultados
	var practicas []models.Practica

	// Construir la consulta con filtros condicionales
	query := database.DB.Model(&models.Practica{})

	if modalidad != "" {
		query = query.Where("modalidad = ?", modalidad)
	}

	if areaPractica != "" {
		query = query.Where("area_practica = ?", areaPractica)
	}
	if jornada != "" {
		query = query.Where("jornada = ?", jornada)
	}

	if ubicacion != "" {
		query = query.Where("ubicacion = ?", ubicacion)
	}

	if fecha != "" {
		query = query.Where("MONTH(fecha_publicacion) = ?", fecha)
	}

	// Ejecutar la consulta y almacenar el resultado en `practicas`
	result := query.Find(&practicas)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las prácticas"})
		return
	}

	// Respuesta exitosa con las prácticas filtradas
	c.JSON(http.StatusOK, gin.H{"practicas": practicas})
}
