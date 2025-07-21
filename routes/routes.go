package routes

import (
	"fieldreserve/constants"
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/service"

	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.IUserController, jwtService service.InterfaceJWTService) {
	routes := route.Group("/api/users")

	routes.POST("/register", userController.CreateUser)
	routes.POST("/login", userController.ReadUserByEmail)

	// Route yang butuh login (token)
	protected := routes.Group("/")
	protected.Use(middleware.Authentication(jwtService))
	protected.PATCH("/update/:id", userController.UpdateUser)
	protected.DELETE("/delete/:id", userController.DeleteUser)

	// Route yang hanya bisa diakses oleh admin
	adminOnly := protected.Group("/")
	adminOnly.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))
	{
		adminOnly.GET("/get-all-user", userController.ReadAllUser)
	}
}
