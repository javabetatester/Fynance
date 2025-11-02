package main

import (
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/infrastructure"
	"Fynance/internal/middleware"
	"Fynance/internal/routes"
	"Fynance/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	db := infrastructure.NewDb()

	userService := user.Service{
		Repository: &infrastructure.UserRepository{DB: db},
	}

	authService := auth.Service{
		Repository: &infrastructure.UserRepository{DB: db},
	}

	goalService := goal.Service{
		Repository: &infrastructure.GoalRepository{DB: db},
	}

	transactionService := transaction.Service{
		Repository:         &infrastructure.TransactionRepository{DB: db},
		CategoryRepository: &infrastructure.TransactionCategoryRepository{DB: db},
	}

	jwtService := utils.NewJwtService()

	handler := routes.Handler{
		UserService:        userService,
		JwtService:         jwtService,
		AuthService:        authService,
		GoalService:        goalService,
		TransactionService: transactionService,
	}

	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/login", handler.Authenticate)
		public.POST("/register", handler.Registration)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(jwtService))
	private.Use(middleware.RequireOwnership())
	{
		private.GET("/users/:id", handler.GetUserByID)
		private.GET("/users/email", handler.GetUserByEmail)
		private.PUT("/users/:id", handler.UpdateUser)
		private.DELETE("/users/:id", handler.DeleteUser)
		private.POST("/goal/create", handler.CreateGoal)
		private.POST("/transaction/create", handler.CreateTransaction)
		private.POST("/transaction/category/create", handler.CreateCategory)
		private.GET("/transaction/list", handler.GetTransactions)
	}
	router.Run(":8080")
}
