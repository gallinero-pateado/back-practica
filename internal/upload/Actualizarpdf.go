package upload

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"practica/internal/storage"
	"strings"

	"github.com/gin-gonic/gin"
)

// UpdatePDFHandler maneja la actualización de archivos PDF de un usuario
// @Summary Actualizar un archivo PDF
// @Description Sube un nuevo archivo PDF a Firebase Storage, elimina el archivo anterior (si existe) y actualiza el campo correspondiente en la base de datos del usuario autenticado
// @Tags upload
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Nuevo archivo PDF a subir"
// @Success 200 {object} map[string]string "URL del PDF actualizado y mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error al actualizar el archivo PDF"
// @Router /update-pdf [post]
func UpdatePDFHandler(c *gin.Context) {
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

	// Verificar que el archivo sea un PDF
	if !isPDF(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo debe ser un PDF"})
		return
	}

	// Buscar al usuario para obtener la URL del archivo PDF anterior
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid).First(&usuario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos del usuario", "details": err.Error()})
		return
	}

	// Eliminar el archivo PDF anterior si existe
	if usuario.CV != "" {
		err = storage.DeleteFileFromFirebase(usuario.CV, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el archivo anterior", "details": err.Error()})
			return
		}
	}

	// Subir el nuevo archivo PDF a Firebase Storage
	url, err := storage.UploadFileToFirebase(file, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir el archivo PDF", "details": err.Error()})
		return
	}

	// Actualizar el campo Cv en la base de datos con la nueva URL
	result := database.DB.Model(&usuario).Where("firebase_usuario = ?", uid).Update("cv", url)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el archivo PDF en la base de datos"})
		return
	}

	// Responder con la URL del archivo PDF actualizado
	c.JSON(http.StatusOK, gin.H{
		"message": "Archivo PDF actualizado correctamente",
		"url":     url,
	})
}

// isPDF verifica que el archivo sea un PDF
func ispff(filename string) bool {
	return strings.ToLower(strings.TrimSpace(filename[len(filename)-4:])) == ".pdf"
}
