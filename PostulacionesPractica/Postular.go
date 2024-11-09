package PostulacionesPractica

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"practica/internal/database"
	"practica/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type PostulacionRequest struct {
	Mensaje string `json:"mensaje"` // Campo para el mensaje del usuario
}

func Postularpractica(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Obtener el ID de la práctica de la URL
	practicaidStr := c.Param("practicaid")
	if practicaidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
		return
	}

	// Convertir practicaidStr a int
	var practicaid int
	if _, err := fmt.Sscan(practicaidStr, &practicaid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
		return
	}

	// Obtener el mensaje del cuerpo de la solicitud
	var req PostulacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar el usuario por el uid de Firebase
	var usuario models.Usuario
	result := database.DB.Where("firebase_usuario = ?", uid).Take(&usuario)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar el usuario en la base de datos"})
		return
	}

	// Guardar el ID del usuario en una variable
	idUsuario := usuario.Id // idUsuario es de tipo uint

	// Verificar si ya existe una postulacion para el usuario y la práctica
	var existingPostulacion models.Postulacion
	err := database.DB.Where("id_usuario = ? AND id_practica = ?", idUsuario, practicaid).First(&existingPostulacion).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Ya has postulado a esta práctica"})
		return
	}

	localTime := time.Now().Local()

	// Si no existe, crear una nueva postulacion
	postulacion := models.Postulacion{
		Id_usuario:            int(idUsuario), // Conversión de uint a int
		Fecha_postulacion:     localTime,
		Id_estado_postulacion: 1,           // Ajusta según tu lógica de estado
		Mensaje:               req.Mensaje, // Mensaje proporcionado por el usuario
		Id_practica:           practicaid,
		Verificacion:          false,
	}
	// Generar un token verificador
	token, err := GenerateVerificationToken(usuario.Correo, uint(postulacion.Id_practica), int(usuario.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de verificación"})
		return
	}

	// Enviar el correo de verificación
	if err := SendVerificationEmail(usuario.Correo, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo de verificación"})
		return
	}
	if err := database.DB.Create(&postulacion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la postulación"})
		return
	}

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Postulación creada exitosamente"})
}

func GenerateVerificationToken(email string, practicaid uint, idusuario int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["email"] = email
	claims["practicaid"] = practicaid
	claims["idusuario"] = idusuario
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // El token expira en 24 horas
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Función para enviar el correo de verificación
func SendVerificationEmail(email, token string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	to := email
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("Subject: Verificación de practica\n\nPor favor verifica tu practica haciendo clic en el siguiente enlace:\n" +
		"http://localhost:8080/verify-email?token=" + token)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)

	return err
}

func VerifyPostulacionHandler(c *gin.Context) {
	tokenString := c.Query("token")

	// Parsear el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Manejar el error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		practicaid := claims["practicaid"].(uint)
		idusuario := claims["idusuario"].(int)

		// Actualizar el estado de validación en la base de datos
		var postulaciones models.Postulacion
		result := database.DB.Model(&postulaciones).Where("id_practica = ?", practicaid).Update("id_usuario", idusuario)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el estado del usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Postulacion verificada exitosamente."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido o expirado"})
	}
}
