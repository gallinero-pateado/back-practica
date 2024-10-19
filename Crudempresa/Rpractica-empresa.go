package Crudempresa

import (
	"fmt"
	"net/http"
	"practica/internal/database"
	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPracticasEmpresas(c *gin.Context) {
	var practicas []models.Practica

	// Obtener el ID de la empresa de la URL
	empresaidN := c.Param("empresaid")
	if empresaidN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica no proporcionado"})
		return
	}

	// Convertir practicaidStr a int
	var empresaid int
	if _, err := fmt.Sscan(empresaidN, &empresaid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de práctica inválido"})
		return
	}

	// Buscar prácticas relacionadas con la empresa en la base de datos
	if err := database.DB.Where("id_empresa = ?", empresaid).Find(&practicas).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prácticas no encontradas"})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, practicas)
}
