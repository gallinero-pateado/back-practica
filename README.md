Backend Célula Prácticas 
Versión 1.2

- Se añadió GORM para el manejo de la base de datos.
- Se añade gin para manejo de solicitudes.
- La función de registro está funcionando correctamente, almacenando en firebase y guardando el usuario en Supabase.

Para ejecutar el codigo debes de tener las credenciales app.env y el json de firebase (solicitarlo)

El código se ejecuta con:
	go run main.go

Para ver la documentación de Swagger:
	http://localhost:8080/swagger/index.html

@xsilvamo @JaimeFigueroaP22 @Xin0tr3

## Configuración para que funcione el archivo Dockerfile

Para evitar incluir las credenciales directamente en la imagen Docker, puedes montar los archivos de configuración necesarios durante la ejecución del contenedor. Sigue estos pasos:

1. **Obtener de los Backend las credenciales y el .json de firebase**:
   - `config/app.env`
   - `config/serviceAccountKey.json`

2. **Ejecución del contenedor con Docker**:
   Usa el siguiente comando para ejecutar el contenedor, asegurándote de reemplazar las rutas de los archivos según tu sistema:

   ```bash
   docker run --name trabajo-b -p 8080:8080 \
     -e GIN_MODE=release \
     -v "ruta/al/archivo/app.env:/root/config/app.env" \
     -v "ruta/al/archivo/serviceAccountKey.json:/root/config/serviceAccountKey.json" \
     back-practica

OBS: Con antelacion debe estar la imagen creada en este caso la imagen se llama "back-practica" y el contenedor "trabajo-b" y esta en el puerto "8080"

@gmigryk @Jeremy2210