package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// FiltroPracticas obtiene las prácticas aplicando filtros opcionales
// @Summary Filtra las prácticas
// @Description Filtra las prácticas según los parámetros opcionales como modalidad, área de práctica, jornada, ubicación y fecha de publicación
// @Tags practicas
// @Accept json
// @Produce json
// @Param modalidad query string false "Filtrar por modalidad (presencial, remoto, etc.)"
// @Param area_practica query string false "Filtrar por área de práctica (ingeniería, marketing, etc.)"
// @Param jornada query string false "Filtrar por jornada (completa, parcial, etc.)"
// @Param ubicacion query string false "Filtrar por ubicación (ciudad, país, etc.)"
// @Param fecha_publicacion query string false "Filtrar por mes de publicación (ejemplo: '2024-10')"
// @Success 200 {object} gin.H "Lista de prácticas filtradas"
// @Failure 500 {object} ErrorResponse "Error al obtener las prácticas"
// @Router /practicas/filtro [get]
// FiltroPracticas aplica filtros opcionales a las prácticas y devuelve el resultado
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
