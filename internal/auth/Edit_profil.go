package auth

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

type EditProfileRequest struct {
	FechaNacimiento *string `json:"fecha_nacimiento,omitempty"`
	AnoIngreso      *string `json:"ano_ingreso,omitempty"`
	IdCarrera       *uint   `json:"id_carrera,omitempty"`
}

// EditProfileHandler permite a los usuarios editar selectivamente su perfil
// @Summary Editar perfil de usuario
// @Description Permite a los usuarios autenticados actualizar selectivamente su perfil
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param profile body EditProfileRequest true "Datos para actualizar el perfil"
// @Success 200 {object} string "Perfil editado correctamente"
// @Failure 400 {object} string "Datos inválidos"
// @Failure 401 {object} string "Usuario no autenticado"
// @Failure 500 {object} string "Error al editar el perfil"
// @Router /edit-profile [patch]
func EditProfileHandler(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener los datos del perfil desde el cuerpo de la solicitud
	var req EditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Crear un mapa con los campos actualizables
	updates := make(map[string]interface{})
	if req.FechaNacimiento != nil {
		updates["fecha_nacimiento"] = *req.FechaNacimiento
	}
	if req.AnoIngreso != nil {
		updates["ano_ingreso"] = *req.AnoIngreso
	}
	if req.IdCarrera != nil {
		updates["id_carrera"] = *req.IdCarrera
	}

	// Verificar si hay campos para actualizar
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay datos para actualizar"})
		return
	}

	// Actualizar el perfil del usuario en la base de datos
	var usuario models.Usuario
	result := database.DB.Model(&usuario).Where("firebase_usuario = ?", uid).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al editar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil editado correctamente"})
}
