package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// practicasRequest estructura de los datos recibidos
type practicasRequest struct {
	Titulo           string    `json:"Titulo"`
	Descripcion      string    `json:"Descripcion"`
	Ubicacion        string    `json:"Ubicacion"`
	Fecha_inicio     time.Time `json:"Fecha_inicio"`
	Fecha_fin        time.Time `json:"Fecha_fin"`
	Requisitos       string    `json:"Requisitos"`
	Fecha_expiracion time.Time `json:"Fecha_expiracion"`
	Modalidad        string    `json:"Modalidad"`
	Area_practica    string    `json:"Area_practica"`
	Jornada          string    `json:"Jornada"`
}

// UpdatePractica actualiza una práctica existente
// @Summary Actualiza una práctica por su ID
// @Description Actualiza los detalles de una práctica existente con los datos proporcionados
// @Tags practicas
// @Accept json
// @Produce json
// @Param id path string true "ID de la práctica a actualizar"
// @Param practica body practicasRequest true "Datos de la práctica actualizada"
// @Success 200 {string} string "La práctica fue actualizada exitosamente"
// @Failure 400 {string} string "Descripción del error de solicitud
// @Failure 404 {string} string "Práctica no encontrada
// @Failure 500 {string} string "Error al actualizar la práctica en la base de datos
// @Router /Upgradepracticas/{id} [put]
func UpdatePractica(c *gin.Context) {
	var req practicasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar la empresa en la base de datos"})
		return
	}

	// Obtener el ID de la práctica a actualizar desde los parámetros de la ruta
	id := c.Param("id")

	var practica models.Practica
	// Buscar la práctica en la base de datos
	if result := database.DB.First(&practica, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Práctica no encontrada"})
		return
	}
	idEmpresa := empresa.Id_empresa
	// Actualizar los campos de la práctica con los datos de la solicitud
	practica.Titulo = req.Titulo
	practica.Descripcion = req.Descripcion
	practica.Id_empresa = int(idEmpresa)
	practica.Ubicacion = req.Ubicacion
	practica.Fecha_inicio = req.Fecha_inicio
	practica.Fecha_fin = req.Fecha_fin
	practica.Requisitos = req.Requisitos
	practica.Fecha_expiracion = req.Fecha_expiracion
	practica.Fecha_publicacion = time.Now().Local() // Hora de actualización
	practica.Modalidad = req.Modalidad
	practica.Area_practica = req.Area_practica
	practica.Jornada = req.Jornada

	// Guardar los cambios en la base de datos
	if result := database.DB.Save(&practica); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la práctica en la base de datos"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "La práctica fue actualizada exitosamente", "id_practica": practica.Id})
}
