package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("🔐 AdminMiddleware: Verificando permisos de admin...")
		
		userRole, exists := c.Get("userRole")
		if !exists {
			log.Println("❌ AdminMiddleware: User role not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			log.Printf("❌ AdminMiddleware: Invalid role type: %T", userRole)
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		log.Printf("🔍 AdminMiddleware: User role: %s", role)
		
		if role != "admin" {
			log.Printf("❌ AdminMiddleware: Access denied. User role: %s, required: admin", role)
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		log.Println("✅ AdminMiddleware: Admin access granted")
		c.Next()
	}
}

