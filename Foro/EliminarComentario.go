package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// EliminarComentario elimina un comentario existente
// @Summary Eliminar un comentario
// @Description Elimina un comentario existente en un tema. El usuario debe estar autenticado y ser el propietario del comentario.
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "ID del comentario"
// @Success 200 {object} map[string]interface{} "Comentario eliminado exitosamente"
// @Failure 400 {object} map[string]interface{} "ID de comentario inválido"
// @Failure 401 {object} map[string]interface{} "Usuario no autenticado"
// @Failure 403 {object} map[string]interface{} "No tienes permiso para eliminar este comentario"
// @Failure 404 {object} map[string]interface{} "Comentario no encontrado"
// @Failure 500 {object} map[string]interface{} "Error al eliminar el comentario"
// @Router /comentarios/{id} [delete]
func EliminarComentario(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, map[string]interface{}{"error": "No tienes permiso para eliminar este comentario"})
		return
	}

	// Eliminar el comentario de la base de datos
	if err := database.DB.Delete(&comentario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al eliminar el comentario"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Comentario eliminado exitosamente"})
}
