package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePractica elimina una práctica por su ID
// @Summary Elimina una práctica por ID
// @Description Elimina una práctica específica de la base de datos utilizando su ID
// @Tags practicas
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "ID de la práctica"
// @Success 200 {string} string "La práctica fue eliminada exitosamente"
// @Failure 400 {string} string "ID inválido"
// @Failure 404 {string} string "Práctica no encontrada"
// @Failure 500 {string} string "Error al eliminar la práctica"
// @Router /Delete-practica/{id} [delete]
func DeletePractica(c *gin.Context) {
	// Obtener el UID de Firebase del contexto
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la práctica desde la URL
	practicaidStr := c.Param("id")
	if practicaidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
		return
	}

	// Convertir practicaidStr a int
	practicaid, err := strconv.Atoi(practicaidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
		return
	}

	// Buscar la empresa en la base de datos asociada al UID de Firebase
	var empresa models.Usuario_empresa
	result := database.DB.Where("firebase_usuario_empresa = ?", uid).First(&empresa)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa en la base de datos"})
		return
	}

	// Guardar el ID de la empresa autenticada
	idEmpresa := empresa.Id_empresa // asumiendo que Id es uint

	// Buscar la práctica en la base de datos
	var practica models.Practica
	result = database.DB.First(&practica, practicaid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Práctica no encontrada"})
		return
	}

	// Verificar si la práctica pertenece a la empresa autenticada
	if practica.Id_empresa != int(idEmpresa) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No autorizado para eliminar esta práctica"})
		return
	}

	// Eliminar la práctica si la empresa autenticada coincide con el dueño de la práctica
	if err := database.DB.Delete(&practica).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la práctica"})
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "La práctica fue eliminada exitosamente"})
}
