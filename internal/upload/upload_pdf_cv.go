package upload

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"practica/internal/storage"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadPDFHandler maneja la subida de archivos PDF y actualiza el campo correspondiente del usuario
// @Summary Subir un archivo PDF
// @Description Sube un archivo PDF a Firebase Storage y actualiza el campo en la base de datos del usuario autenticado
// @Tags upload
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Archivo PDF a subir"
// @Success 200 {object} map[string]string "URL del PDF subido y mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 401 {object} map[string]string "Usuario no autenticado"
// @Failure 500 {object} map[string]string "Error al subir el archivo PDF"
// @Router /upload-pdf [post]
func UploadPDFHandler(c *gin.Context) {
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

	// Subir el archivo PDF a Firebase Storage
	url, err := storage.UploadFileToFirebase(file, "ulink-sprint-1.appspot.com") // Reemplaza con tu bucket
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir el archivo PDF", "details": err.Error()})
		return
	}

	// Actualizar el campo en la base de datos con la URL del PDF
	var usuario models.Usuario
	result := database.DB.Model(&usuario).Where("firebase_usuario = ?", uid).Update("cv", url)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el campo PDF en la base de datos"})
		return
	}

	// Responder con la URL del archivo PDF subido
	c.JSON(http.StatusOK, gin.H{
		"message": "Archivo PDF subido y campo actualizado correctamente",
		"url":     url,
	})
}

// isPDF verifica que el archivo sea un PDF
func isPDF(filename string) bool {
	return strings.ToLower(strings.TrimSpace(filename[len(filename)-4:])) == ".pdf"
}
