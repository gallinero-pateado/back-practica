package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Leer todos los comentarios de un tema específico
func LeerComentarios(c *gin.Context) {
	temaIDStr := c.Param("id")                         // Obtener el ID del tema desde la URL
	temaID, err := strconv.ParseInt(temaIDStr, 10, 64) // Convertir el ID del tema a int64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tema inválido"})
		return
	}

	var comentarios []models.Comentario
	if err := database.DB.Where("tema_id = ?", temaID).Find(&comentarios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los comentarios"})
		return
	}

	if len(comentarios) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron comentarios para este tema"})
		return
	}

	c.JSON(http.StatusOK, comentarios) // Devuelve la lista de comentarios
}
