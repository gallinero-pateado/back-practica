package Crudempresa

import (
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// practicaRequest representa la estructura de los datos recibidos para crear una nueva práctica
type practicaRequest struct {
	Id               int       `gorm:"primaryKey;autoIncrement"`
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

// ErrorResponse representa la estructura de un error
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse representa una respuesta exitosa con el ID de la práctica
type SuccessResponse struct {
	Message    string `json:"message"`
	IdPractica int    `json:"id_practica"`
}

// Createpractica crea una nueva oferta de práctica
// @Summary Crea una nueva oferta de práctica
// @Description Crea una oferta de práctica basada en los datos proporcionados y la guarda en la base de datos
// @Tags practicas
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param practica body practicaRequest true "Datos de la nueva práctica"
// @Success 200 {object} SuccessResponse "Oferta de práctica creada exitosamente con el ID de la práctica"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 500 {object} ErrorResponse "Error al guardar la práctica en la base de datos"
// @Router /Create-practicas [post]
func Createpractica(c *gin.Context) {
	var req practicaRequest
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
	idEmpresa := empresa.Id_empresa

	// Crear práctica
	practica := models.Practica{
		Titulo:             req.Titulo,
		Descripcion:        req.Descripcion,
		Id_empresa:         int(idEmpresa), // Usar el ID de la empresa encontrada
		Ubicacion:          req.Ubicacion,
		Fecha_inicio:       req.Fecha_inicio,
		Fecha_fin:          req.Fecha_fin,
		Requisitos:         req.Requisitos,
		Fecha_expiracion:   req.Fecha_expiracion,
		Fecha_publicacion:  time.Now().Local(), // Hora de creación
		Modalidad:          req.Modalidad,
		Area_practica:      req.Area_practica,
		Jornada:            req.Jornada,
		Id_estado_practica: 1, // Suponiendo que 1 es el estado inicial
	}

	// Guardar la práctica en la base de datos
	result = database.DB.Create(&practica)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la práctica en la base de datos"})
		return
	}

	// Respuesta exitosa con el ID de la práctica creada
	c.JSON(http.StatusOK, SuccessResponse{
		Message:    "La oferta de práctica fue creada exitosamente",
		IdPractica: int(practica.Id), // Convertir de uint a int
	})
}
