package notificaciones

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MensajeNotificacion struct {
	Mensaje string `json:"mensaje"`
}

// Notificacion Manda una notificacion
// @Summary Enviar una notificación
// @Description Enviar una notificación con un mensaje específico
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param mensaje body MensajeNotificacion true "Mensaje de la notificación"
// @Success 200 {object} gin.H{"mensaje": "string"}
// @Failure 400 {object} gin.H{"error": "Datos inválidos"}
// @Router /notificacion [post]
func ProcesarNotificacion(c *gin.Context, mensaje string) {
	var req MensajeNotificacion

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	//websocket.Broadcast <- websocket.MensajeNotificacion{Mensaje: req.Mensaje}
	c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})
}

func Notificacion(c *gin.Context) {
	mensaje := "Mensaje de prueba"
	ProcesarNotificacion(c, mensaje)
}
