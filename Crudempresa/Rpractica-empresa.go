package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// GetPracticasEmpresas obtiene las prácticas asociadas a la empresa del usuario autenticado
// @Summary Obtiene las prácticas de la empresa del usuario autenticado
// @Description Recupera las prácticas asociadas a la empresa del usuario autenticado mediante su UID
// @Tags practicas
// @Accept json
// @Produce json
// @Success 200 {array} models.Practica "Lista de prácticas"
// @Failure 401 {object} gin.H{"error": "Usuario no autenticado"}
// @Failure 404 {object} gin.H{"error": "Prácticas no encontradas"}
// @Router /GetPracticasEmpresas [get]
func GetPracticasEmpresas(c *gin.Context) {
	var practicas []models.Practica

	// Obtener el UID de Firebase del contexto
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar la empresa en la base de datos asociada al UID de Firebase
	var empresa models.Usuario_empresa
	result := database.DB.Where("firebase_usuario_empresa = ?", uid).First(&empresa)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
		return
	}

	// Buscar prácticas relacionadas con la empresa en la base de datos
	if err := database.DB.Where("id_empresa = ?", empresa.Id_empresa).Find(&practicas).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prácticas no encontradas"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, practicas)
}
