package PostulacionesPractica

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

// PostulacionResponse representa la estructura de una postulación en la respuesta.
type PostulacionResponse struct {
	PracticaID       uint   `json:"practica_id"`
	NombreEstudiante string `json:"nombre_estudiante"`
	Correo           string `json:"correo"`
	Mensaje          string `json:"mensaje"`
	FechaPostulacion string `json:"fecha_postulacion"`
}

// ErrorResponse representa un mensaje de error en la respuesta.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ObtenerPostulantesPorPractica obtiene la lista de estudiantes que postularon a una práctica específica.
//
// @Summary Obtiene los postulantes a una práctica específica
// @Description Devuelve la lista de estudiantes que postularon a una práctica específica perteneciente a la empresa autenticada.
// @Tags Postulaciones
// @Param Authorization header string true "Bearer <token>"
// @Param practicaid path string true "ID de la práctica"
// @Success 200 {object} []PostulacionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /practicas/{practicaid}/postulaciones [get]
func ObtenerPostulantesPorPractica(c *gin.Context) {
	// Obtener el token de autenticación
	empresaUID, exists := c.Get("uid") // Suponiendo que el UID de la empresa está incluido en el contexto
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Empresa no autenticada"})
		return
	}

	// Buscar la empresa asociada al UID
	var empresa models.Usuario_empresa
	if err := database.DB.Where("firebase_usuario_empresa = ?", empresaUID).First(&empresa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la empresa"})
		return
	}

	// Obtener el ID de la práctica de la URL
	practicaID := c.Param("practicaid")
	if practicaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
		return
	}

	// Verificar si la práctica pertenece a la empresa autenticada
	var practica models.Practica
	if err := database.DB.Where("id = ? AND id_empresa = ?", practicaID, empresa.Id_empresa).First(&practica).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Práctica no encontrada o no pertenece a esta empresa"})
		return
	}

	// Buscar los postulantes asociados a esta práctica
	var postulantes []models.Postulacion
	if err := database.DB.Where("id_practica = ?", practicaID).Find(&postulantes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener postulaciones"})
		return
	}

	// Crear un slice para almacenar las postulaciones
	var postulaciones []struct {
		PracticaID       uint   `json:"practica_id"`
		NombreEstudiante string `json:"nombre_estudiante"`
		Correo           string `json:"correo"`
		Mensaje          string `json:"mensaje"`
		FechaPostulacion string `json:"fecha_postulacion"`
	}

	// Iterar sobre las postulaciones y agregar la información del estudiante
	for _, postulante := range postulantes {
		var usuario models.Usuario
		if err := database.DB.Where("id = ?", postulante.Id_usuario).First(&usuario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del estudiante"})
			return
		}

		postulaciones = append(postulaciones, struct {
			PracticaID       uint   `json:"practica_id"`
			NombreEstudiante string `json:"nombre_estudiante"`
			Correo           string `json:"correo"`
			Mensaje          string `json:"mensaje"`
			FechaPostulacion string `json:"fecha_postulacion"`
		}{
			PracticaID:       practica.Id,
			NombreEstudiante: usuario.Nombres,
			Correo:           usuario.Correo,
			Mensaje:          postulante.Mensaje,
			FechaPostulacion: postulante.Fecha_postulacion.Format("2006-01-02 15:04:05"),
		})
	}

	// Retornar la lista de postulaciones
	c.JSON(http.StatusOK, gin.H{"postulaciones": postulaciones})
}
