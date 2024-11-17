package notificaciones

import (
	"net/http"
	"practica/internal/database"
	models "practica/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

func Create_notificacion(c *gin.Context) {
	var req models.Notificaciones_All
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Notificacion := models.Notificaciones_All{
		Id:                 req.Id,
		Mensaje:            req.Mensaje,
		Fecha_hora_mensaje: time.Now().Local(),
		Estado:             "Enviado",
	}

	if err := database.DB.Create(&Notificacion).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
