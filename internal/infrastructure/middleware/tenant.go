package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener tenant ID del header
		tenantID := c.GetHeader("X-Tenant-ID")

		var tenantIDInt int
		var err error

		// Si no viene en el header, intentar obtenerlo del token JWT
		if tenantID == "" {
			if tokenTenantID, exists := c.Get("tenantID"); exists {
				tenantIDInt = tokenTenantID.(int)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header required"})
				c.Abort()
				return
			}
		} else {
			tenantIDInt, err = strconv.Atoi(tenantID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
				c.Abort()
				return
			}
		}

		// Guardar tenant ID en el contexto
		c.Set("tenantID", tenantIDInt)
		c.Next()
	}
}

