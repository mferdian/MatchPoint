package middleware

import (
	"context"
	"fieldreserve/constants"
	"fieldreserve/service"
	"fieldreserve/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(jwtService service.InterfaceJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userID, err := jwtService.GetUserIDByToken(tokenStr)
		if err != nil {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		role, err := jwtService.GetRoleByToken(tokenStr)
		if err != nil {
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set("Authorization", tokenStr)
		ctx.Set("user_id", userID)
		ctx.Set("role", role)

		// Simpan token ke context.Context agar bisa diakses di service layer
		stdCtx := context.WithValue(ctx.Request.Context(), "token", tokenStr)
		ctx.Request = ctx.Request.WithContext(stdCtx)

		ctx.Next()
	}
}
