package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Leer todos los comentarios de un tema específico
// @Summary Obtener los comentarios de un tema
// @Description Devuelve todos los comentarios de un tema específico dado su ID. Si no hay comentarios, devuelve un mensaje indicando que no se encontraron comentarios.
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "ID del tema"
// @Success 200 {array} models.Comentario "Lista de comentarios"
// @Failure 400 {object} map[string]interface{} "ID de tema inválido"
// @Failure 404 {object} map[string]interface{} "No se encontraron comentarios para este tema"
// @Failure 500 {object} map[string]interface{} "Error al obtener los comentarios"
// @Router /temas/{id}/comentarios [get]
func LeerComentarios(c *gin.Context) {
	temaIDStr := c.Param("id")                         // Obtener el ID del tema desde la URL
	temaID, err := strconv.ParseInt(temaIDStr, 10, 64) // Convertir el ID del tema a int64
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID de tema inválido"})
		return
	}

	var comentarios []models.Comentario
	if err := database.DB.Where("tema_id = ?", temaID).Find(&comentarios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al obtener los comentarios"})
		return
	}

	if len(comentarios) == 0 {
		c.JSON(http.StatusNotFound, map[string]interface{}{"message": "No se encontraron comentarios para este tema"})
		return
	}

	c.JSON(http.StatusOK, comentarios) // Devuelve la lista de comentarios
}
