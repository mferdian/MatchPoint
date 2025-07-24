package routes

import (
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userController controller.IUserController, categoryController controller.ICategoryController, jwtService service.InterfaceJWTService) {
	user := r.Group("/api/users")
	user.Use(middleware.Authentication(jwtService))

	user.PATCH("/update-profile/:id", userController.UpdateUser)
	user.GET("/get-detail-user/:id", userController.GetUserByID)
	user.DELETE("/delete-profile/:id", userController.DeleteUser)
	user.GET("/get-all-categories", categoryController.GetAllCatgory)
}
