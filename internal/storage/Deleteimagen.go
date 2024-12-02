package storage

import (
	"context"
	"fmt"
	"strings"
)

// DeleteFileFromFirebase elimina un archivo desde Firebase Storage utilizando su URL
func DeleteFileFromFirebase(fileURL string, bucketName string) error {
	// Inicializar Firebase Storage
	app, err := InitFirebaseStorage()
	if err != nil {
		return err
	}

	ctx := context.Background()
	client, err := app.Storage(ctx)
	if err != nil {
		return fmt.Errorf("error obteniendo cliente de storage: %v", err)
	}

	// Obtener el bucket de Firebase Storage
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("error obteniendo bucket: %v", err)
	}

	// Extraer el nombre del archivo desde la URL
	// La URL de Firebase Storage tiene el siguiente formato:
	// https://firebasestorage.googleapis.com/v0/b/{bucket}/o/{filename}?alt=media
	// Necesitamos extraer el nombre del archivo después de "/o/"
	// Usamos strings.Split para separar la URL en partes
	parts := strings.Split(fileURL, "/o/")
	if len(parts) < 2 {
		return fmt.Errorf("URL no válida para eliminar el archivo: %v", fileURL)
	}
	fileName := parts[1]

	// Eliminar el archivo del bucket
	object := bucket.Object(fileName)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("error al eliminar el archivo: %v", err)
	}

	return nil
}
