package main

import (
	"fieldreserve/cmd"
	"fieldreserve/config/database"
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/repository"
	"fieldreserve/routes"
	"fieldreserve/service"
	"fieldreserve/utils" // tambahkan ini
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ==== Set up logger ====
	utils.SetUpLogger()
	utils.Log.Info("Logger initialized")

	// ==== Load .env ====
	if err := godotenv.Load(); err != nil {
		utils.Log.Warn("No .env file found")
	}

	// ==== Database ====
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// ==== Seeder command ====
	if len(os.Args) > 1 {
		utils.Log.WithField("command", os.Args[1]).Info("Seeder command triggered")
		cmd.Command(db)
		return
	}

	// ==== Inisialisasi ====
	var (
		jwtService = service.NewJWTService()

		userRepo       = repository.NewUserRepository(db)
		userService    = service.NewUserService(userRepo, jwtService)
		userController = controller.NewUserController(userService)

		categoryRepo       = repository.NewCategoryRepository(db)
		categoryService    = service.NewCategoryService(categoryRepo)
		categoryController = controller.NewCategoryController(categoryService)

		fieldRepo       = repository.NewFieldRepository(db)
		fieldService    = service.NewFieldService(fieldRepo)
		fieldController = controller.NewFieldController(fieldService)

		scheduleRepo     = repository.NewScheduleRepository(db)
		scheduleService  = service.NewScheduleService(scheduleRepo, fieldRepo)
		scheduleController = controller.NewScheduleController(scheduleService)

		bookingRepo       = repository.NewBookingRepository(db)
		bookingService    = service.NewBookingService(bookingRepo, jwtService, scheduleRepo, fieldRepo)
		bookingController = controller.NewBookingController(bookingService)
	)

	// ==== Router ====
	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.PublicRoutes(server, userController)
	routes.UserRoutes(server, userController, categoryController, fieldController, scheduleController, bookingController, jwtService)
	routes.AdminRoutes(server, userController, categoryController, fieldController, scheduleController, bookingController, jwtService)

	server.Static("/assets", "./assets")

	// ==== Port ====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	utils.Log.WithField("address", serve).Info("Starting server...")

	if err := server.Run(serve); err != nil {
		utils.Log.WithError(err).Fatal("Error running server")
	}
}
