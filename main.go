package main

import (
	"fieldreserve/cmd"
	"fieldreserve/config/database"
	"fieldreserve/controller"
	"fieldreserve/middleware"
	"fieldreserve/repository"
	"fieldreserve/routes"
	"fieldreserve/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Seeder command
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

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

		scheduleRepo = repository.NewScheduleRepository(db)
		scheduleService = service.NewScheduleService(scheduleRepo, fieldRepo)
		scheduleController = controller.NewScheduleController(scheduleService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.PublicRoutes(server, userController)
	routes.UserRoutes(server, userController, categoryController, fieldController, scheduleController, jwtService)
	routes.AdminRoutes(server, userController, categoryController, fieldController, scheduleController, jwtService)

	server.Static("/assets", "./assets")

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

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
