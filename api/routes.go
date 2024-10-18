package api

import (
	Cempresa "practica/Crudempresa"
	postular "practica/Postulacionespractica"
	"practica/internal/auth"
	"practica/internal/upload"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Cambia el puerto si es necesario
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	router.POST("/register", auth.RegisterHandler)
	router.POST("/login", auth.LoginHandler)
	router.POST("/register_empresa", auth.RegisterHandler_empresa)
	router.GET("/verify-email", auth.VerifyEmailHandler)
	router.POST("/password-reset", auth.SendPasswordResetEmailHandler)
	router.POST("/resend-verification", auth.ResendVerificationEmailHandler)
	router.GET("/usuarios/:uid", auth.GetUsuariouid)
	router.GET("/get-allusuario", auth.GetAllUsuarios)

	// rutas crud practicas
	router.POST("/Create-practicas", Cempresa.Createpractica)
	router.GET("/Get-practicas", Cempresa.GetAllPracticas)
	router.PUT("/Update-practicas/:id", Cempresa.UpdatePractica)
	router.DELETE("/Delete-practica/:id", Cempresa.DeletePractica)

	//filtros pagina
	router.GET("/filtro-practicas", Cempresa.FiltroPracticas)

	// Rutas protegidas
	protected := router.Group("/").Use(auth.AuthMiddleware) // Agrupar las rutas protegidas con el middleware
	{
		protected.POST("/complete-profile", auth.CompleteProfileHandler) // Ruta para completar perfil
		protected.POST("/upload-image", upload.UploadImageHandler)       // Ruta para subir imágenes
		protected.GET("/profile-status", auth.GetProfileStatusHandler)   // Ruta para obtener el estado del perfil
		protected.POST("/postulacion-practicas/:practicaid", postular.Postularpractica)
	}

	return router
}
