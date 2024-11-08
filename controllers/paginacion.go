package controllers

import (
	"net/http"
	"practica/internal/database"
	"strconv"

	"practica/internal/models"

	"github.com/gin-gonic/gin"
)

func ObtenerComentarios(c *gin.Context) {
	var comentarios []models.Comentario
	var total int64

	temaIDStr := c.Param("id") // Obtener el ID del tema desde la URL
	temaID, err := strconv.ParseInt(temaIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tema inválido"})
		return
	}
	// Convertir el ID del tema a int64

	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Calcular el offset
	offset := (page - 1) * limit

	// Obtener el total de comentarios
	database.DB.Model(&models.Comentario{}).Where("tema_id = ?", temaID).Count(&total)

	// Obtener los comentarios paginados
	database.DB.Offset(offset).Limit(limit).Find(&comentarios)

	// Devolver la respuesta con paginación
	c.JSON(http.StatusOK, gin.H{
		"data":       comentarios,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}
