package notificaciones

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MensajeNotificacion struct {
	Mensaje string `json:"mensaje"`
}

// ProcesarNotificacion procesa una notificación y devuelve un mensaje
// @Summary Procesar notificación
// @Description Esta función procesa la notificación y devuelve un mensaje especificado.
// @Tags notificaciones
// @Accept json
// @Produce json
// @Param mensaje body string true "Mensaje de la notificación"
// @Success 200 {object} map[string]string "Mensaje procesado"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Router /notificaciones [post]
func ProcesarNotificacion(c *gin.Context, mensaje string) {
	var req MensajeNotificacion

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})
}

// Notificacion maneja la solicitud de notificación y procesa el mensaje de prueba
// @Summary Enviar notificación de prueba
// @Description Devuelve un mensaje de prueba como respuesta
// @Tags notificaciones
// @Success 200 {object} map[string]string "Mensaje de prueba enviado correctamente"
// @Router /notificaciones/test [get]
func Notificacion(c *gin.Context) {
	mensaje := "Mensaje de prueba"
	ProcesarNotificacion(c, mensaje)
}
