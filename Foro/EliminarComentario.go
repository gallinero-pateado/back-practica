package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// EliminarComentario elimina un comentario existente
func EliminarComentario(c *gin.Context) {
	comentarioID := c.Param("id") // Obtener el ID del comentario desde la URL
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar el comentario en la base de datos
	var comentario models.Comentario
	if err := database.DB.First(&comentario, comentarioID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentario no encontrado"})
		return
	}

	// Buscar el ID del usuario en la base de datos usando el UID
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid.(string)).First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar que el comentario pertenece al usuario autenticado
	if comentario.UsuarioID != usuario.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para eliminar este comentario"})
		return
	}

	// Eliminar el comentario de la base de datos
	if err := database.DB.Delete(&comentario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el comentario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comentario eliminado exitosamente"})
}
