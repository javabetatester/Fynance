package main

import (
	"Fynance/internal/domain/login"
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
	loginService := login.Service{
		Repository: &infrastructure.LoginRepository{DB: db},
	}

	jwtService := utils.NewJwtService()

	handler := routes.Handler{
		LoginService: loginService,
		UserService:  userService,
	}

	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/login", handler.Login)
		public.POST("/users", handler.CreateUser)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(jwtService))
	{
		private.GET("/users/:id", handler.GetUserByID)
		private.GET("/users/email/:email", handler.GetUserByEmail)
		private.PUT("/users/:id", handler.UpdateUser)
		private.DELETE("/users/:id", handler.DeleteUser)
	}

	router.Run(":8080")
}
