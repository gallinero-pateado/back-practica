package main

import (
	//"fmt"
	"log"

	//"net/http"
	"practica/api"
	"practica/internal/auth"
	"practica/internal/database"
	"practica/internal/models"
	"practica/internal/storage"
	"practica/pkg/config"
	"practica/websocket"

	_ "practica/docs" // Importar la documentación generada

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Cargar las configuraciones desde el archivo .env
	err := config.LoadConfig("config/app.env")
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar la base de datos
	err = database.InitDatabase()
	if err != nil {
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}

	// Realizar la migración de la tabla fuera de la transacción
	if !database.DB.Migrator().HasTable(&models.Usuario{}) {
		err = database.DB.AutoMigrate(&models.Usuario{})
		if err != nil {
			log.Fatalf("Error al migrar modelos: %v", err)
		}
	}

	// Realizar la migración de la tabla fuera de la transacción
	if !database.DB.Migrator().HasTable(&models.Usuario_empresa{}) {
		err = database.DB.AutoMigrate(&models.Usuario_empresa{})
		if err != nil {
			log.Fatalf("Error al migrar modelos: %v", err)
		}
	}

	// Inicializar Firebase
	err = auth.InitFirebase()
	if err != nil {
		log.Fatalf("Error inicializando Firebase: %v", err)
	}

	// Inicializar Firebase Storage
	_, err = storage.InitFirebaseStorage()
	if err != nil {
		log.Fatalf("Error inicializando Firebase Storage: %v", err)
	}

	//Registrar rutas
	router := api.SetupRoutes()

	// Agregar la ruta de Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/*go func() {
		for {
			msg := <-websocket.Broadcast
			websocket.Mutex.Lock()
			for _, cliente := range websocket.Clientes {
				err := cliente.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("error sending message: %v", err)
					cliente.Conn.Close()
					delete(websocket.Clientes, fmt.Sprintf("%d", cliente.Id))
				}
			}
			websocket.Mutex.Unlock()
		}
	}()*/

	go websocket.HandleMessages()
	// Mostrar en consola que el WebSocket está activo
	log.Println("WebSocket server is active and listening for connections")
	// Iniciar el servidor
	router.Run(":8080")
}
