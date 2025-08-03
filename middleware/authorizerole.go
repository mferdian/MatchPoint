package middleware

import (
	"fieldreserve/constants"
	"fieldreserve/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, ok := c.Get("role")
		if !ok {
			utils.Log.Warn("Authorization failed: role not found in context")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, "unauthorized: role not found", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		roleStr, ok := roleVal.(string)
		if !ok {
			utils.Log.Warn("Authorization failed: role is not a valid string")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, "unauthorized: invalid role format", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(roleStr, allowed) {
				utils.Log.Infof("Authorized access for role: %s", roleStr)
				c.Next()
				return
			}
		}

		utils.Log.Warnf("Forbidden access for role: %s", roleStr)
		res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, "forbidden: you don't have access to this resource", nil)
		c.AbortWithStatusJSON(http.StatusForbidden, res)
	}
}
