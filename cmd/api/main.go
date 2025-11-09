package main

import (
	"Fynance/config"
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/infrastructure"
	"Fynance/internal/logger"
	"Fynance/internal/middleware"
	"Fynance/internal/routes"

	docs "Fynance/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Fynance API
// @version 1.0
// @description API de gestão financeira pessoal (Fynance)
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Utilize "Bearer <token>" no header Authorization

func main() {
	cfg := config.Load()
	logger.Init(cfg)

	db, err := infrastructure.NewDb(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Falha ao inicializar banco de dados")
	}

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
		Repository:  goalRepo,
		UserService: userService,
	}

	transactionService := transaction.Service{
		Repository:         transactionRepo,
		CategoryRepository: categoryRepo,
		UserService:        &userService,
	}

	investmentService := investment.Service{
		Repository:      investmentRepo,
		TransactionRepo: transactionRepo,
		UserService:     &userService,
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

	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api")
	{
		public.POST("/auth/login", handler.Authenticate)
		public.POST("/auth/register", handler.Registration)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(jwtService))
	private.Use(middleware.RequireOwnership())
	{

		goals := private.Group("/goals")
		{
			goals.POST("", handler.CreateGoal)
			goals.PATCH("/:id", handler.UpdateGoal)
			goals.GET("", handler.ListGoals)
			goals.GET("/:id", handler.GetGoal)
			goals.DELETE("/:id", handler.DeleteGoal)
		}

		transactions := private.Group("/transactions")
		{
			transactions.POST("", handler.CreateTransaction)
			transactions.GET("", handler.GetTransactions)
			transactions.GET("/:id", handler.GetTransaction)
			transactions.PATCH("/:id", handler.UpdateTransaction)
			transactions.DELETE("/:id", handler.DeleteTransaction)
		}

		categories := private.Group("/categories")
		{
			categories.POST("", handler.CreateCategory)
			categories.GET("", handler.ListCategories)
			categories.PATCH("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
		}

		investments := private.Group("/investments")
		{
			investments.POST("", handler.CreateInvestment)
			investments.GET("", handler.ListInvestments)
			investments.GET("/:id", handler.GetInvestment)
			investments.POST("/:id/contribution", handler.MakeContribution)
			investments.POST("/:id/withdraw", handler.MakeWithdraw)
			investments.GET("/:id/return", handler.GetInvestmentReturn)
			investments.DELETE("/:id", handler.DeleteInvestment)
			investments.PATCH("/:id", handler.UpdateInvestment)
		}
	}

	serverAddr := ":" + cfg.Server.Port
	logger.Info().
		Str("address", serverAddr).
		Str("environment", cfg.App.Environment).
		Msg("Servidor iniciando")

	if err := router.Run(serverAddr); err != nil {
		logger.Fatal().Err(err).Msg("Falha ao iniciar servidor")
	}
}
