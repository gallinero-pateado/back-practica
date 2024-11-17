package subcanales

import (
	"net/http"
	"practica/internal/database"
	models "practica/internal/models"

	"github.com/gin-gonic/gin"
)

func Create_canal(c *gin.Context) {
	var req models.Subcanales
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Canales := models.Subcanales{
		Id:       req.Id,
		Id_canal: req.Id_canal,
		Nombre:   req.Nombre,
	}

	if err := database.DB.Create(&Canales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
