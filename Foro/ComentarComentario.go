package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv" // Importar para convertir el ID del tema y el ID del comentario padre

	"github.com/gin-gonic/gin"
)

// Añadir un comentario o respuesta a un tema
func ReposnderComentario(c *gin.Context) {
	var nuevoComentario models.Comentario
	if err := c.ShouldBindJSON(&nuevoComentario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	temaIDStr := c.Param("id")                         // Obtener el ID del tema desde la URL
	temaID, err := strconv.ParseInt(temaIDStr, 10, 64) // Convertir el ID del tema a int64
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tema inválido"})
		return
	}

	// Obtener el ID del comentario padre desde la URL, si existe
	comentarioPadreIDStr := c.DefaultQuery("comentario_padre_id", "") // Opcional: Se pasa como parámetro en la URL
	var comentarioPadreID *uint
	if comentarioPadreIDStr != "" {
		comentarioPadreIDParsed, err := strconv.ParseUint(comentarioPadreIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de comentario padre inválido"})
			return
		}
		comentarioPadreID = new(uint)
		*comentarioPadreID = uint(comentarioPadreIDParsed)
	}

	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar el ID del usuario en la base de datos usando el UID
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid.(string)).First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	nuevoComentario.TemaID = temaID                       // Asignar el ID del tema
	nuevoComentario.UsuarioID = usuario.Id                // Asignar el ID del usuario autenticado
	nuevoComentario.ComentarioPadreID = comentarioPadreID // Asignar el ID del comentario padre si existe

	if err := database.DB.Create(&nuevoComentario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el comentario"})
		return
	}

	c.JSON(http.StatusCreated, nuevoComentario)
}
