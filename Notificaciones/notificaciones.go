package notificaciones

import (
	"fmt"
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	websocket "practica/websocket"
	"time"

	"github.com/gin-gonic/gin"
)

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

func GuardarNotificacion(idMensaje int, idReceptor int) error {
	NotificacionTodos := models.Notificaciones_All{
		Fecha_hora_mensaje: time.Now(),
		Estado:             "Enviado",
	}
	return database.DB.Create(&NotificacionTodos).Error
}

func NotificarNuevoHilo(c *gin.Context) {
	var Comentarios []models.Comentario
	var Notificaciones []models.Notificaciones_All
	if err := database.DB.Where("nuevo_comentario = ?", true).Find(&Comentarios).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentario no encontrado"})
		return
	}

	var Tema models.Tema
	for _, comentario := range Comentarios {
		if err := database.DB.Where("id = ?", comentario.TemaID).First(&Tema).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var seguidores []models.Usuario
	for _, comentario := range Comentarios {
		// Obtener los usuarios que siguen el hilo
		if err := database.DB.Where("comentario_padre_id = ?", comentario.Comentario_padre_id).Find(&seguidores).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, comentario := range Comentarios {
		for _, usuario := range seguidores {
			notificacion := models.Notificaciones_All{
				Id:                 usuario.Id,
				Titulo:             "Nuevo comentario",
				Mensaje:            fmt.Sprintf("Nuevo comentario en el hilo que sigues: %s", Tema.Titulo),
				Fecha_hora_mensaje: time.Now(),
				Estado:             "Enviado",
			}
			Notificaciones = append(Notificaciones, notificacion)
			websocket.EnviarMensajeATemaOComentario(Tema.Titulo, fmt.Sprint(comentario.Comentario_padre_id), notificacion)
		}
	}

	for _, comentario := range Comentarios {
		comentario.Nuevo_comentario = false
		if err := database.DB.Save(&comentario).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Guardar las notificaciones en la base de datos
	if err := database.DB.Create(&Notificaciones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificaciones enviadas"})
}
