package main

import (
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/infrastructure"
	"Fynance/internal/middleware"
	"Fynance/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := infrastructure.NewDb()

	userRepo := &infrastructure.UserRepository{DB: db}
	goalRepo := &infrastructure.GoalRepository{DB: db}
	transactionRepo := &infrastructure.TransactionRepository{DB: db}
	categoryRepo := &infrastructure.TransactionCategoryRepository{DB: db}
	investmentRepo := &infrastructure.InvestmentRepository{DB: db}

	userService := user.Service{
		Repository: userRepo,
	}

	authService := auth.Service{
		Repository:  userRepo,
		UserService: &userService,
	}

	goalService := goal.Service{
		Repository: goalRepo,
	}

	transactionService := transaction.Service{
		Repository:         transactionRepo,
		CategoryRepository: categoryRepo,
	}

	investmentService := investment.Service{
		Repository:      investmentRepo,
		TransactionRepo: transactionRepo,
	}

	jwtService := middleware.NewJwtService(&userService)

	handler := routes.Handler{
		UserService:        userService,
		JwtService:         jwtService,
		AuthService:        authService,
		GoalService:        goalService,
		TransactionService: transactionService,
		InvestmentService:  investmentService,
	}

	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/auth/login", handler.Authenticate)
		public.POST("/auth/register", handler.Registration)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(jwtService))
	private.Use(middleware.RequireOwnership())
	{
		users := private.Group("/users")
		{
			users.GET("/:id", handler.GetUserByID)
			users.GET("/email", handler.GetUserByEmail)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		goals := private.Group("/goals")
		{
			goals.POST("", handler.CreateGoal)
		}

		transactions := private.Group("/transactions")
		{
			transactions.POST("", handler.CreateTransaction)
			transactions.GET("", handler.GetTransactions)
		}

		categories := private.Group("/categories")
		{
			categories.POST("", handler.CreateCategory)
		}

		investments := private.Group("/investments")
		{
			investments.POST("", handler.CreateInvestment)
			investments.GET("", handler.ListInvestments)
			investments.GET("/:id", handler.GetInvestment)
			investments.POST("/:id/contribution", handler.MakeContribution)
			investments.POST("/:id/withdrawal", handler.MakeWithdrawal)
			investments.GET("/:id/return", handler.GetInvestmentReturn)
			investments.DELETE("/:id", handler.DeleteInvestment)
		}
	}

	router.Run(":8080")
}
