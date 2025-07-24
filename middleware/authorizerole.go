package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "unauthorized: role not found",
			})
			return
		}

		roleStr := userRole.(string)
		for _, allowed := range allowedRoles {
			if strings.EqualFold(roleStr, allowed) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "forbidden: you don't have access to this resource",
		})
	}
}
