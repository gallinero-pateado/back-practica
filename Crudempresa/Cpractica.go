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
	Titulo             string    `json:"Titulo"`
	Descripcion        string    `json:"Descripcion"`
	Id_empresa         int       `json:"Id_Empresa"`
	Ubicacion          string    `json:"Ubicacion"`
	Fecha_inicio       time.Time `json:"Fecha_inicio"`
	Fecha_fin          time.Time `json:"Fecha_fin"`
	Requisitos         string    `json:"Requisitos"`
	Fecha_expiracion   time.Time `json:"Fecha_expiracion"`
	Id_estado_practica int       `json:"Id_estado_practica"`
	Fecha_publicacion  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Modalidad          string    `json:"Modalidad"`
	Area_practica      string    `json:"Area_practica"`
	Jornada            string    `json:"Jornada"`
}

// Createpractica crea una nueva oferta de práctica
// @Summary Crea una nueva oferta de práctica
// @Description Crea una oferta de práctica basada en los datos proporcionados y la guarda en la base de datos
// @Tags practicas
// @Accept json
// @Produce json
// @Param practica body practicaRequest true "Datos de la nueva práctica"
// @Success 200 {object} gin.H "Oferta de práctica creada exitosamente con el ID de la práctica"
// @Failure 400 {object} ErrorResponse "Datos inválidos"
// @Failure 500 {object} ErrorResponse "Error al guardar la práctica en la base de datos"
// @Router /Createpracticas [post]
// Createpractica maneja la creación de una nueva oferta de práctica
func Createpractica(c *gin.Context) {
	var req practicaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tomar la hora de creación
	localTime := time.Now().Local()

	// Crear práctica
	practica := models.Practica{
		Titulo:             req.Titulo,
		Descripcion:        req.Descripcion,
		Id_empresa:         req.Id_empresa,
		Ubicacion:          req.Ubicacion,
		Fecha_inicio:       req.Fecha_inicio,
		Fecha_fin:          req.Fecha_fin,
		Requisitos:         req.Requisitos,
		Fecha_expiracion:   req.Fecha_expiracion,
		Fecha_publicacion:  localTime,
		Modalidad:          req.Modalidad,
		Area_practica:      req.Area_practica,
		Jornada:            req.Jornada,
		Id_estado_practica: 1,
	}

	// Guardar la práctica en la base de datos
	result := database.DB.Create(&practica)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la practica en la base de datos"})
		return
	}

	// Respuesta exitosa con el ID de la práctica creada
	c.JSON(http.StatusOK, gin.H{"message": "La Oferta de practica fue creada exitosamente", "id_practica": practica.Id})
}
