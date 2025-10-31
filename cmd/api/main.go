package main

import (
	login "Fynance/internal/domain/auth"
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
		Repository: &infrastructure.UserRepository{DB: db},
	}

	jwtService := utils.NewJwtService()

	handler := routes.Handler{
		UserService:  userService,
		JwtService:   jwtService,
		LoginService: loginService,
	}

	router := gin.Default()

	public := router.Group("/api")
	{
		public.POST("/login", handler.Login)
		public.POST("/users", handler.CreateUser)
	}

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(jwtService))
	private.Use(middleware.RequireOwnership())
	{
		private.GET("/users/:id", handler.GetUserByID)
		private.GET("/users/email", handler.GetUserByEmail)
		private.PUT("/users/:id", handler.UpdateUser)
		private.DELETE("/users/:id", handler.DeleteUser)
	}

	router.Run(":8080")
}
