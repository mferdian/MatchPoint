package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)


func GetUserIDAndRoleFromContext(c *gin.Context) (string, string, error) {
    userID, ok := c.Get("user_id")
    if !ok {
        return "", "", errors.New("user_id not found in context")
    }

    role, ok := c.Get("Role")
    if !ok {
        return "", "", errors.New("role not found in context")
    }

    return userID.(string), role.(string), nil
}


