package database

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	//websocket "practica/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckPostulacionForChanges checks for changes in the postulacion table
// @Summary Verificar cambios en la tabla de postulaciones
// @Description Verifica cambios en la tabla de postulaciones y envía correos electrónicos si hay cambios
// @Tags postulaciones
// @Produce json
// @Success 200 {string} string "Verificación completada"
// @Failure 500 {object} ErrorResponse "Error al verificar cambios"
// @Router /check-postulaciones [get]
func CheckPostulacionForChanges(c *gin.Context, db *gorm.DB) {
	for {
		var postulaciones []Postulacion

		// Query to check for changes in the postulacion table
		if err := db.Table("postulacion").
			Select(`postulacion.*, "Usuario".correo, estado_postulacion.nom_estado_postulacion`).
			Joins(`left join "Usuario" on postulacion.id_usuario = "Usuario".id`).
			Joins("left join estado_postulacion on postulacion.id_estado_postulacion = estado_postulacion.id_estado_postulacion").
			Where("postulacion.id_estado_postulacion IN (2, 3) AND postulacion.previo_estado_postulacion = 1").
			Find(&postulaciones).Error; err != nil {
			log.Fatalf("Error querying the database: %v", err)
		}

		for _, postulacion := range postulaciones {
			// Send email to the user
			err := SendEmail(postulacion.Correo, postulacion.NomEstadoPostulacion)
			if err != nil {
				log.Printf("Error sending email to %s: %v", postulacion.Correo, err)
			}
			//clienteID := postulacion.IDUsuario
			//msg := "Subject: Cambio en el estado de tu postulación\n\nEl estado de tu postulación ha cambiado a: " + postulacion.NomEstadoPostulacion

			// Actualiza el previo_estado_postulacion al actual estado_postulacion
			if err := db.Model(&Postulacion{}).
				Where("id = ?", postulacion.ID).
				Update("previo_estado_postulacion", postulacion.IDEstadoPostulacion).Error; err != nil {
				log.Printf("Error updating previo_estado_postulacion for postulacion ID %d: %v", postulacion.ID, err)
			}
		}
	}
}

func CambiarEstadoPostulacion(c *gin.Context, db *gorm.DB) {
	for {
		var postulaciones []Postulacion

		usuarioIdStr := c.Param("usuarioId")
		if usuarioIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
			return
		}

		// Convertir practicaidStr a int
		var usuarioId int
		if _, err := fmt.Sscan(usuarioIdStr, &usuarioId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
			return
		}

		eleccionStr := c.Param("eleccionId")
		if eleccionStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
			return
		}

		// Convertir practicaidStr a int
		var eleccionId int
		if _, err := fmt.Sscan(eleccionStr, &eleccionId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
			return
		}

		postulacionidStr := c.Param("postulacionid")
		if postulacionidStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
			return
		}

		// Convertir practicaidStr a int
		var postulacionid int
		if _, err := fmt.Sscan(postulacionidStr, &postulacionid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
			return
		}

		// Chequea si existe una postulacion para el usuario y la practica
		if err := db.Table("postulacion").
			Select(`postulacion.*, "Usuario".correo, estado_postulacion.nom_estado_postulacion`).
			Joins(`left join "Usuario" on postulacion.id_usuario = "Usuario".id`).
			Joins("left join estado_postulacion on postulacion.id_estado_postulacion = estado_postulacion.id_estado_postulacion").
			Where("postulacion.id_usuario = ? AND postulacion.id = ?", usuarioId, postulacionid).
			Find(&postulaciones).Error; err != nil {
			log.Fatalf("Error querying the database: %v", err)
		}

		for _, postulacion := range postulaciones {
			//Guarda el estado anterior
			if err := db.Model(&Postulacion{}).
				Where("id = ?", postulacion.ID).
				Update("previo_estado_postulacion", postulacion.IDEstadoPostulacion).Error; err != nil {
				log.Printf("Error updating previo_estado_postulacion for postulacion ID %d: %v", postulacion.ID, err)
			}

			if eleccionId == 2 {
				postulacion.IDEstadoPostulacion = 2
			} else if eleccionId == 3 {
				postulacion.IDEstadoPostulacion = 3
			}
			// Actualiza el estado_postulacion al estado actualizado
			if err := db.Model(&Postulacion{}).
				Where("id = ?", postulacion.ID).
				Update("estado_postulacion", postulacion.IDEstadoPostulacion).Error; err != nil {
				log.Printf("Error actualizando previo_estado_postulacion en postulacion ID %d: %v", postulacion.ID, err)
			}

			// Send email to the user
			err := SendEmail(postulacion.Correo, postulacion.NomEstadoPostulacion)
			if err != nil {
				log.Printf("Error sending email to %s: %v", postulacion.Correo, err)
			}
			//clienteID := postulacion.IDUsuario
			//msg := "Subject: Cambio en el estado de tu postulación\n\nEl estado de tu postulación ha cambiado a: " + postulacion.NomEstadoPostulacion

			// Actualiza el previo_estado_postulacion al actual estado_postulacion
		}
	}
}

// SendEmail sends an email to the user
// @Summary Send an email notification
// @Description Sends an email to the user with the updated status of their application
// @Tags email
// @Accept json
// @Produce json
// @Param to query string true "Recipient email address"
// @Param estadoPostulacion query string true "Application status"
// @Success 200 {string} string "Email sent successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /sendemail [post]
func SendEmail(to string, estadoPostulacion string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("Subject: Cambio en el estado de tu postulación\n\nEl estado de tu postulación ha cambiado a: " + estadoPostulacion)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)

	return err
}

// CheckPostulacionForChangesHandler maneja la ruta para verificar cambios en la tabla de postulaciones
// @Summary Verificar cambios en la tabla de postulaciones
// @Description Verifica cambios en la tabla de postulaciones y envía correos electrónicos si hay cambios
// @Tags postulaciones
// @Produce json
// @Success 200 {string} string "Verificación completada"
// @Failure 500 {object} ErrorResponse "Error al verificar cambios"
// @Router /check-postulaciones [get]
func CheckPostulacionForChangesHandler(c *gin.Context) {
	go CheckPostulacionForChanges(c, DB)
	c.JSON(http.StatusOK, gin.H{"message": "Verificación iniciada"})
}

// HandleSendEmail handles the route for sending emails
func HandleSendEmail(c *gin.Context) {
	to := c.Query("to")
	estadoPostulacion := c.Query("estadoPostulacion")

	if to == "" || estadoPostulacion == "" {
		c.JSON(400, gin.H{"error": "Datos de solicitud inválidos"})
		return
	}

	err := SendEmail(to, estadoPostulacion)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al enviar el correo"})
		return
	}
	c.JSON(200, gin.H{"message": "Correo enviado con éxito"})
}

// ErrorResponse estructura de la respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// CheckNuevoPostulanteForChanges checks for changes in the postulacion table
// @Summary Verificar cambios en la tabla de postulaciones
// @Description Verifica cambios en la tabla de postulaciones y envía correos electrónicos si hay cambios
// @Tags postulaciones
// @Produce json
// @Success 200 {string} string "Verificación completada"
// @Failure 500 {object} ErrorResponse "Error al verificar cambios"
// @Router /Check-NuevoPostulanteForChanges [get]
func CheckNuevoPostulanteForChanges(db *gorm.DB) {
	mensaje := "Tiene una nueva postulacion en su practica"

	for {
		var postulaciones []Postulacion

		// Query to check for changes in the postulacion table
		if err := db.Table("postulacion").
			Select(`postulacion.*, "Usuario_empresa".correo_empresa, postulacion.nueva_postulacion`).
			Joins(`left join practica on postulacion.id_practica = practica.id`).
			Joins(`left join "Usuario_empresa" on practica.id_empresa = "Usuario_empresa".id_empresa`).
			Where("postulacion.nueva_postulacion = ?", true).
			Find(&postulaciones).Error; err != nil {
			log.Fatalf("Error querying the database: %v", err)
		}

		for _, postulacion := range postulaciones {
			// Manda un correo a la empresa
			MandarCorreoNuevoPostulante(postulacion.CorreoEmpresa, mensaje)

			// Actualiza el previo_estado_postulacion al actual estado_postulacion
			if err := db.Model(&Postulacion{}).
				Where("id = ?", postulacion.ID).
				Update("nueva_postulacion", false).Error; err != nil {
				log.Printf("Error updating previo_estado_postulacion for postulacion ID %d: %v", postulacion.ID, err)
			}
		}
	}
}

// MandarCorreoNuevoPostulante sends an email to the user
// @Summary Send an email notification
// @Description Sends an email to the user with the updated status of their application
// @Tags email
// @Accept json
// @Produce json
// @Param to query string true "Recipient email address"
// @Param estadoPostulacion query string true "Application status"
// @Success 200 {string} string "Email sent successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /MandarCorreoNuevoPostulante [post]
func MandarCorreoNuevoPostulante(to string, mensaje string) {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte(mensaje)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
	}
}

func HandleMandarCorreoNuevoPostulante(c *gin.Context) {
	to := c.Query("to")
	mensaje := c.Query("mensaje")

	if to == "" {
		c.JSON(400, gin.H{"error": "Datos de solicitud inválidos"})
		return
	}

	MandarCorreoNuevoPostulante(to, mensaje)
	c.JSON(200, gin.H{"message": "Correo enviado con éxito"})
}

// CheckNuevoPostulanteForChangesHandler maneja la ruta para verificar cambios en la tabla de postulaciones
// @Summary Verificar cambios en la tabla de postulaciones
// @Description Verifica cambios en la tabla de postulaciones y envía correos electrónicos si hay cambios
// @Tags postulaciones
// @Produce json
// @Success 200 {string} string "Verificación completada"
// @Failure 500 {object} ErrorResponse "Error al verificar cambios"
// @Router /Check-NuevoPostulanteForChanges [get]
func CheckNuevoPostulanteForChangesHandler(c *gin.Context) {
	go CheckNuevoPostulanteForChanges(DB)
	c.JSON(http.StatusOK, gin.H{"message": "Verificación iniciada"})
}

type Postulacion struct {
	ID                      uint      `json:"id"`
	IDUsuario               uint      `json:"id_usuario"`
	IDEmpresa               uint      `json:"id_empresa"`
	IDPractica              uint      `json:"id_practica"`
	FechaPostulacion        time.Time `json:"fecha_postulacion"`
	Mensaje                 string    `json:"mensaje"`
	IDEstadoPostulacion     uint      `json:"id_estado_postulacion"`
	PrevioEstadoPostulacion uint      `json:"previo_estado_postulacion"`
	Correo                  string    `json:"correo"`
	CorreoEmpresa           string    `json:"correo_empresa"`
	NomEstadoPostulacion    string    `json:"nom_estado_postulacion"`
	NuevaPostulacion        bool      `json:"nueva_postulacion"`
}

func (Postulacion) TableName() string {
	return "postulacion"
}
