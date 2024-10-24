package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware verifica si el usuario tiene uno de los roles permitidos
func RoleMiddleware(rolesPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el rol del usuario desde el contexto
		Rol, exists := c.Get("rol")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se proporcionó el rol del usuario"})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		for _, rolPermitido := range rolesPermitidos {
			if Rol == rolPermitido {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para acceder a esta ruta"})
		c.Abort()
	}
}
