package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// ActualizarComentario actualiza un comentario existente
// @Summary Actualiza un comentario
// @Description Actualiza un comentario específico en la base de datos si el usuario tiene permiso para hacerlo
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del comentario a actualizar"
// @Param comentario body models.Comentario true "Datos del comentario actualizado"
// @Success 200 {object} models.Comentario "Comentario actualizado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 401 {object} map[string]interface{} "Usuario no autenticado"
// @Failure 403 {object} map[string]interface{} "No tienes permiso para actualizar este comentario"
// @Failure 404 {object} map[string]interface{} "Comentario no encontrado"
// @Failure 500 {object} map[string]interface{} "Error al actualizar el comentario"
// @Router /comentarios/{id} [put]
func ActualizarComentario(c *gin.Context) {
	var comentarioActualizado models.Comentario
	if err := c.ShouldBindJSON(&comentarioActualizado); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	comentarioID := c.Param("id") // Obtener el ID del comentario desde la URL
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Usuario no autenticado"})
		return
	}

	// Buscar el comentario en la base de datos
	var comentario models.Comentario
	if err := database.DB.First(&comentario, comentarioID).Error; err != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Comentario no encontrado"})
		return
	}

	// Buscar el ID del usuario en la base de datos usando el UID
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid.(string)).First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Usuario no encontrado"})
		return
	}

	// Verificar que el comentario pertenece al usuario autenticado
	if comentario.UsuarioID != usuario.Id {
		c.JSON(http.StatusForbidden, map[string]interface{}{"error": "No tienes permiso para actualizar este comentario"})
		return
	}

	// Actualizar el comentario con los nuevos datos
	if err := database.DB.Model(&comentario).Updates(comentarioActualizado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al actualizar el comentario"})
		return
	}

	c.JSON(http.StatusOK, comentario)
}
