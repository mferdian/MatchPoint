package routes

import (
	"fieldreserve/constants"
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/service"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine, userController controller.IUserController, categoryController controller.ICategoryController, fieldcontroller controller.IFieldController, scheduleController controller.IScheduleController, bookingController controller.IBookingController,
	jwtService service.InterfaceJWTService) {
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

	// Field Management
	admin.POST("/create-field", fieldcontroller.CreateField)
	admin.PATCH("/update-field/:id", fieldcontroller.UpdateField)
	admin.DELETE("/delete-field/:id", fieldcontroller.DeleteField)

	// Schedule Management
	admin.POST("/create-schedule", scheduleController.CreateSchedule)
	admin.PATCH("/update-schedule/:id", scheduleController.UpdateSchedule)
	admin.DELETE("/delete-schedule/:id", scheduleController.DeleteSchedule)
	admin.GET("/get-all-schedule", scheduleController.GetAllSchedule)
}
