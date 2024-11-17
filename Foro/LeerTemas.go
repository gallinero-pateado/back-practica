package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// Leer todos los temas
// @Summary Obtener todos los temas
// @Description Devuelve una lista de todos los temas disponibles. Si no hay temas, devuelve un mensaje indicando que no se encontraron temas.
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Tema "Lista de temas"
// @Failure 404 {object} map[string]interface{} "No se encontraron temas"
// @Failure 500 {object} map[string]interface{} "Error al obtener los temas"
// @Router /temas [get]
func LeerTemas(c *gin.Context) {
	var temas []models.Tema
	if err := database.DB.Find(&temas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al obtener los temas"})
		return
	}

	if len(temas) == 0 {
		c.JSON(http.StatusNotFound, map[string]interface{}{"message": "No se encontraron temas"})
		return
	}

	c.JSON(http.StatusOK, temas) // Devuelve la lista de temas
}
