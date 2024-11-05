package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// ActualizarComentario actualiza un comentario existente
func ActualizarComentario(c *gin.Context) {
	var comentarioActualizado models.Comentario
	if err := c.ShouldBindJSON(&comentarioActualizado); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para actualizar este comentario"})
		return
	}

	// Actualizar el comentario con los nuevos datos
	if err := database.DB.Model(&comentario).Updates(comentarioActualizado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el comentario"})
		return
	}

	c.JSON(http.StatusOK, comentario)
}
