package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// Leer todos los temas
func LeerTemas(c *gin.Context) {
	var temas []models.Tema
	if err := database.DB.Find(&temas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los temas"})
		return
	}

	if len(temas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron temas"})
		return
	}

	c.JSON(http.StatusOK, temas) // Devuelve la lista de temas
}
