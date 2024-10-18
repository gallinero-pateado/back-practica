package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePractica elimina una práctica por su ID
// @Summary Elimina una práctica por ID
// @Description Elimina una práctica específica de la base de datos utilizando su ID
// @Tags practicas
// @Accept json
// @Produce json
// @Param id path int true "ID de la práctica"
// @Success 200 {object} gin.H "La práctica fue eliminada exitosamente"
// @Failure 400 {object} ErrorResponse "ID inválido"
// @Failure 404 {object} ErrorResponse "Práctica no encontrada"
// @Failure 500 {object} ErrorResponse "Error al eliminar la práctica"
// @Router /Deletepracticas/{id} [delete]
// DeletePractica elimina una práctica de la base de datos usando el ID proporcionado en la URL
func DeletePractica(c *gin.Context) {
	// Obtener el ID de la URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convertir el ID de string a int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar la práctica en la base de datos
	var practica models.Practica
	result := database.DB.First(&practica, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Práctica no encontrada"})
		return
	}

	// Eliminar la práctica
	if err := database.DB.Delete(&practica).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la práctica"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "La práctica fue eliminada exitosamente"})
}
