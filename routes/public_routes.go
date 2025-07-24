package routes

import (
	"fieldreserve/controller"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.Engine, userController controller.IUserController) {
	public := r.Group("/api/users")
	public.POST("/register", userController.CreateUser)
	public.POST("/login", userController.GetUserByEmail)
}
