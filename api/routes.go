package api

import (
	Cempresa "practica/Crudempresa"
	foro "practica/Foro"
	notificaciones "practica/Notificaciones"

	postular "practica/PostulacionesPractica"
	controllers "practica/controllers"
	"practica/internal/auth"
	"practica/internal/database"
	"practica/internal/upload"
	websocket "practica/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Cambia el puerto si es necesario
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
	router.GET("/Get-practicas", Cempresa.GetAllPracticas)
	//filtros pagina
	router.GET("/filtro-practicas", Cempresa.FiltroPracticas)
	//leer comentarios
	router.GET("/temas/:id/comentarios", foro.LeerComentarios)
	//leer temas
	router.GET("/temas", foro.LeerTemas)

	router.GET("/obtener-comentarios", controllers.ObtenerComentarios)

	// Rutas protegidas
	protected := router.Group("/").Use(auth.AuthMiddleware) // Agrupar las rutas protegidas con el middleware
	{
		protected.POST("/complete-profile", auth.CompleteProfileHandler)                // Ruta para completar perfil
		protected.POST("/upload-image", upload.UploadImageHandler)                      // Ruta para subir imágenes
		protected.GET("/profile-status", auth.GetProfileStatusHandler)                  // Ruta para obtener el estado del perfil
		protected.POST("/postulacion-practicas/:practicaid", postular.Postularpractica) // Ruta para postular a practicas como usuario
		protected.DELETE("/Delete-practica/:id", Cempresa.DeletePractica)               // Ruta para borrar practica como empresa
		protected.POST("/Create-practicas", Cempresa.Createpractica)                    // Ruta para Crear Practicas como empresa
		protected.GET("/Get-practicas-empresa", Cempresa.GetPracticasEmpresas)          //Ruta para que la empresa vea sus practicas
		protected.PUT("/Update-practicas/:id", Cempresa.UpdatePractica)                 //Ruta para Cambiar datos de practica
		protected.POST("/temas", foro.CrearTema)                                        // Crear un nuevo tema
		protected.POST("/temas/:id/comentarios", foro.AñadirComentario)                 // Añadir un comentario a un tema
		protected.PUT("/comentarios/:id", foro.ActualizarComentario)                    // Actualizar un comentario
		protected.DELETE("/comentarios/:id", foro.EliminarComentario)                   // Eliminar un comentario
	}

	// Rutas de correos
	router.POST("/sendEmail", database.HandleSendEmail)

	// Chequea cambios
	router.GET("/check-postulaciones", database.CheckPostulacionForChangesHandler)
	router.GET("/Check-NuevoPostulanteForChanges", database.CheckNuevoPostulanteForChangesHandler)

	router.GET("/ws/:id", websocket.Handle_WebSocket())
	router.POST("/notificar/:id", notificaciones.NotificarNuevoHilo)
	return router
}
