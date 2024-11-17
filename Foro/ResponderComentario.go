package foro

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"strconv" // Importar para convertir el ID del tema y el ID del comentario padre

	"github.com/gin-gonic/gin"
)

// ResponderComentario permite añadir un comentario o respuesta a un tema
// @Summary Añadir un comentario o respuesta a un tema
// @Description Permite a un usuario añadir un comentario o responder a un comentario en un tema. Se puede incluir un comentario padre (respuesta).
// @Tags Foro
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "ID del tema"
// @Param comentario_padre_id query int false "ID del comentario padre (opcional)"
// @Success 201 {object} models.Comentario "Comentario creado exitosamente"
// @Failure 400 {object} map[string]interface{} "Error en la solicitud (por ejemplo, ID de tema o comentario inválido)"
// @Failure 401 {object} map[string]interface{} "Usuario no autenticado"
// @Failure 500 {object} map[string]interface{} "Error al crear el comentario"
// @Router /comentarios/{id}/respuesta [post]
func ResponderComentario(c *gin.Context) {
	var nuevoComentario models.Comentario
	if err := c.ShouldBindJSON(&nuevoComentario); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	temaIDStr := c.Param("id")                         // Obtener el ID del tema desde la URL
	temaID, err := strconv.ParseInt(temaIDStr, 10, 64) // Convertir el ID del tema a int64
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID de tema inválido"})
		return
	}

	// Obtener el ID del comentario padre desde la URL, si existe
	comentarioPadreIDStr := c.DefaultQuery("comentario_padre_id", "") // Opcional: Se pasa como parámetro en la URL
	var comentarioPadreID *uint
	if comentarioPadreIDStr != "" {
		comentarioPadreIDParsed, err := strconv.ParseUint(comentarioPadreIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID de comentario padre inválido"})
			return
		}
		comentarioPadreID = new(uint)
		*comentarioPadreID = uint(comentarioPadreIDParsed)
	}

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

	nuevoComentario.TemaID = temaID                       // Asignar el ID del tema
	nuevoComentario.UsuarioID = usuario.Id                // Asignar el ID del usuario autenticado
	nuevoComentario.ComentarioPadreID = comentarioPadreID // Asignar el ID del comentario padre si existe

	if err := database.DB.Create(&nuevoComentario).Error; err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error al crear el comentario"})
		return
	}

	c.JSON(http.StatusCreated, nuevoComentario)
}
