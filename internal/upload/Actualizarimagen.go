package upload

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"practica/internal/storage"

	"github.com/gin-gonic/gin"
)

// UpdateImageHandler maneja la actualización de imágenes de perfil del usuario
// @Summary Actualizar una imagen de perfil
// @Description Sube una nueva imagen a Firebase Storage, elimina la anterior (si existe) y actualiza el perfil del usuario autenticado
// @Tags upload
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Nueva imagen a subir"
// @Success 200 {object} map[string]string "URL de la imagen actualizada y mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error al actualizar la imagen"
// @Router /update-image [post]
func UpdateImageHandler(c *gin.Context) {
	// Obtener el UID del usuario autenticado desde el contexto
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se ha proporcionado un archivo"})
		return
	}

	// Buscar al usuario en la base de datos para obtener la URL de la imagen anterior
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid).First(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos del usuario", "details": err.Error()})
		return
	}

	// Si el usuario tiene una foto de perfil, eliminarla
	if usuario.Foto_perfil != "" {
		err := storage.DeleteFileFromFirebase(usuario.Foto_perfil, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la imagen anterior", "details": err.Error()})
			return
		}
	}

	// Subir la nueva imagen a Firebase Storage
	url, err := storage.UploadFileToFirebase(file, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir la nueva imagen", "details": err.Error()})
		return
	}

	// Actualizar la base de datos con la nueva URL de la imagen
	if err := database.DB.Model(&usuario).Where("firebase_usuario = ?", uid).Update("foto_perfil", url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la foto de perfil en la base de datos", "details": err.Error()})
		return
	}

	// Responder con la URL de la nueva imagen
	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen actualizada correctamente",
		"url":     url,
	})
}
