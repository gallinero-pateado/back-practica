package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// CrearTema crea un nuevo tema en el foro
// @Summary Crear un nuevo tema
// @Description Crea un nuevo tema en el foro. El usuario debe estar autenticado para poder crear un tema.
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tema body models.Tema true "Datos del nuevo tema"
// @Success 201 {object} models.Tema "Tema creado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 401 {object} map[string]interface{} "Usuario no autenticado"
// @Failure 500 {object} map[string]interface{} "Error al crear el tema"
// @Router /temas [post]
func CrearTema(c *gin.Context) {
	var nuevoTema models.Tema
	if err := c.ShouldBindJSON(&nuevoTema); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	// Obtener UID del usuario autenticado desde el contexto
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Usuario no autenticado"})
		return
	}

	// Buscar el ID del usuario en la base de datos usando el UID
	var usuario models.Usuario
	if err := database.DB.Where("firebase_usuario = ?", uid.(string)).First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Usuario no encontrado"})
		return
	}

	nuevoTema.UsuarioID = usuario.Id // Asignar el ID del usuario autenticado

	// Crear el nuevo tema
	if err := database.DB.Create(&nuevoTema).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al crear el tema"})
		return
	}

	c.JSON(http.StatusCreated, nuevoTema)
}
