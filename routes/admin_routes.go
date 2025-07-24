package routes

import (
	"fieldreserve/constants"
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/service"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine, userController controller.IUserController, categoryController controller.ICategoryController, jwtService service.InterfaceJWTService) {
	admin := r.Group("/api/admin")
	admin.Use(middleware.Authentication(jwtService))
	admin.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))

	// User management
	admin.GET("/get-all-users", userController.GetAllUser)

	// Category management
	admin.GET("/get-category/:id", categoryController.GetCategoryByID)
	admin.POST("/create-category", categoryController.CreateCategory)
	admin.PATCH("/update-category/:id", categoryController.UpdateCategory)
	admin.DELETE("/detele-category/:id", categoryController.DeleteCategory)
}
