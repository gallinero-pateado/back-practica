package auth

import (
	"context"
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// RegisterRequest estructura de los datos recibidos
type RegisterRequest_empresa struct {
	Email_empresa  string `json:"Email_empresa" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Nombre_empresa string `json:"Nombre_empresa" binding:"required"`
}

// RegisterResponse estructura de la respuesta de registro
type RegisterResponse_empresa struct {
	Message     string `json:"message"`
	FirebaseUID string `json:"firebase_uid"`
}

// RegisterHandler_empresa maneja el registro del usuario
// @Summary Registra un nuevo usuario
// @Description Crea un nuevo usuario en Firebase y lo guarda en la base de datos local
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest_empresa true "Datos del usuario a registrar"
// @Success 200 {object} RegisterResponse_empresa "Usuario registrado correctamente"
// @Failure 400 {object} RegisterResponse_empresa "Solicitud inválida"
// @Failure 500 {object} RegisterResponse_empresa "Error interno del servidor"
// @Router /register_empresa [post]
// RegisterHandler maneja el registro del usuario
func RegisterHandler_empresa(c *gin.Context) {
	var req RegisterRequest_empresa
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crear el usuario en Firebase con email y password
	params := (&auth.UserToCreate{}).
		Email(req.Email_empresa).
		Password(req.Password)

	user, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario en Firebase: " + err.Error()})
		return
	}

	// Crear el usuario en la base de datos sin almacenar la contraseña
	usuario_empresa := models.Usuario_empresa{
		Correo_empresa:           req.Email_empresa,
		Nombre_empresa:           req.Nombre_empresa,
		Firebase_usuario_empresa: user.UID,
		Perfil_Completado:        false,
		Rol:                      "empresa", // Rol por defecto
	}

	result := database.DB.Create(&usuario_empresa)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el usuario en la base de datos"})
		return
	}

	// Generar token de verificación de correo
	token, err := GenerateVerificationToken(req.Email_empresa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de verificación"})
		return
	}

	// Enviar correo de verificación
	err = SendVerificationEmail(req.Email_empresa, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar el correo de verificación"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Usuario empresa creado correctamente. Verifica tu correo", "firebase_uid": user.UID})
}
