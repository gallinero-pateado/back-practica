package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// GetAllPracticas obtiene todas las prácticas
// @Summary Obtiene todas las prácticas
// @Description Recupera la lista completa de todas las prácticas registradas en la base de datos
// @Tags practicas
// @Accept json
// @Produce json
// @Success 200 {array} models.Practica "Lista de todas las prácticas"
// @Failure 500 {string} string "Error al obtener las prácticas"
// @Router /Get-practicas [get]
// GetAllPracticas devuelve todas las prácticas almacenadas en la base de datos
func GetAllPracticas(c *gin.Context) {
	var practicas []models.Practica

	// Obtener todas las prácticas de la base de datos
	if err := database.DB.Find(&practicas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las prácticas"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, practicas)
}
