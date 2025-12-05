package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener tenant ID del token JWT (debe estar disponible desde AuthMiddleware)
		tokenTenantID, exists := c.Get("tenantID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found in token"})
			c.Abort()
			return
		}

		tokenTenantIDInt := tokenTenantID.(int)

		// Obtener tenant ID del header
		headerTenantID := c.GetHeader("X-Tenant-ID")

		var headerTenantIDInt int

		// Si viene en el header, validar que coincida con el del token
		if headerTenantID != "" {
			var err error
			headerTenantIDInt, err = strconv.Atoi(headerTenantID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
				c.Abort()
				return
			}

			// Validar que el tenant del header coincida con el del token
			// Esto previene que un usuario acceda a datos de otro tenant
			if headerTenantIDInt != tokenTenantIDInt {
				c.JSON(http.StatusForbidden, gin.H{"error": "Tenant ID in header does not match token tenant ID"})
				c.Abort()
				return
			}
		}

		// Usar el tenant ID del token (que es la fuente de verdad)
		// Guardar tenant ID en el contexto
		c.Set("tenantID", tokenTenantIDInt)
		c.Next()
	}
}

