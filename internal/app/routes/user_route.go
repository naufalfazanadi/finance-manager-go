package routes

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserRoutes handles user-related routes and initializes dependencies per module
func UserRoutes(api fiber.Router, db *gorm.DB, validator *validator.Validator) {
	// Initialize dependencies for user module
	userRepo := repositories.NewUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, validator)
	userHandler := handlers.NewUserHandler(userUseCase)

	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
