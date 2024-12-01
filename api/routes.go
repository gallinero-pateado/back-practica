package api

import (
	Cempresa "practica/Crudempresa"
	foro "practica/Foro"
	postular "practica/PostulacionesPractica"
	"practica/internal/auth"
	"practica/internal/database"
	"practica/internal/upload"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "https://practicas.tssw.info", "https://descuentos.tssw.info", "https://roomies.tssw.info"}, // Cambia el puerto si es necesario
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

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

	// Rutas protegidas
	protected := router.Group("/").Use(auth.AuthMiddleware) // Agrupar las rutas protegidas con el middleware
	{
		protected.POST("/complete-profile", auth.CompleteProfileHandler)                // Ruta para completar perfil
		protected.POST("/complete-profile/empresa", auth.CompleteProfileEmpresaHandler) // Ruta para completar perfil
		protected.POST("/upload-image", upload.UploadImageHandler)                      // Ruta para subir imágenes
		protected.POST("/update-image", upload.UpdateImageHandler)                      //Ruta Actualizar imagen
		protected.POST("/upload-pdf", upload.UploadPDFHandler)                          //Ruta para subir CV
		protected.POST("/update-pdf", upload.UploadPDFHandler)                          //Ruta actualizar CV
		protected.GET("/profile-status", auth.GetProfileStatusHandler)                  // Ruta para obtener el estado del perfil
		protected.GET("/profile-status/empresa", auth.GetProfileStatusEmpresaHandler)   // Ruta para obtener el estado del perfil de empresa
		protected.POST("/postulacion-practicas/:practicaid", postular.Postularpractica) // Ruta para postular a practicas como usuario
		protected.DELETE("/Delete-practica/:id", Cempresa.DeletePractica)               // Ruta para borrar practica como empresa
		protected.POST("/Create-practicas", Cempresa.Createpractica)                    // Ruta para Crear Practicas como empresa
		protected.GET("/Get-practicas-empresa", Cempresa.GetPracticasEmpresas)          //Ruta para que la empresa vea sus practicas
		protected.PUT("/Update-practicas/:id", Cempresa.UpdatePractica)                 //Ruta para Cambiar datos de practica
		protected.POST("/temas", foro.CrearTema)                                        // Crear un nuevo tema
		protected.POST("/temas/:id/comentarios", foro.AñadirComentario)                 // Añadir un comentario a un tema
		protected.PUT("/comentarios/:id", foro.ActualizarComentario)                    // Actualizar un comentario
		protected.DELETE("/comentarios/:id", foro.EliminarComentario)                   // Eliminar un comentario
		protected.POST("/comentarios/:id/respuesta", foro.ResponderComentario)          //esto es para responder comentarios
		protected.PATCH("/edit-profile", auth.EditProfileHandler)                       // Ruta para editar el perfil
		protected.GET("/practicas/:practicaid/postulaciones", postular.ObtenerPostulantesPorPractica)
	}

	// Rutas de correos
	router.POST("/sendEmail", database.HandleSendEmail)

	router.GET("/check-postulaciones", database.CheckPostulacionForChangesHandler)
	router.GET("/Check-NuevoPostulanteForChanges", database.CheckNuevoPostulanteForChangesHandler)

	return router
}
