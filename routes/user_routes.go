package routes

import (
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	r *gin.Engine,
	userController controller.IUserController,
	categoryController controller.ICategoryController,
	fieldController controller.IFieldController,
	jwtService service.InterfaceJWTService,
) {
	user := r.Group("/api/users")
	user.Use(middleware.Authentication(jwtService))

	// --- User Routes ---
	user.PATCH("/update-profile/:id", userController.UpdateUser)
	user.GET("/get-detail-user/:id", userController.GetUserByID)
	user.DELETE("/delete-profile/:id", userController.DeleteUser)

	// --- Category Routes ---
	user.GET("/get-all-categories", categoryController.GetAllCatgory)

	// --- Field Routes ---
	user.GET("/get-all-field", fieldController.GetAllField)
	user.GET("/get-detail-field/:id", fieldController.GetFieldByID)
}
