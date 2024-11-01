package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// Crear un nuevo tema
func CrearTema(c *gin.Context) {
	var nuevoTema models.Tema
	if err := c.ShouldBindJSON(&nuevoTema); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener UID del usuario autenticado desde el contexto
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

	nuevoTema.UsuarioID = usuario.Id // Asignar el ID del usuario autenticado

	// Crear el nuevo tema
	if err := database.DB.Create(&nuevoTema).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el tema"})
		return
	}

	c.JSON(http.StatusCreated, nuevoTema)
}
