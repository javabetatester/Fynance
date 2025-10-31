package routes

import (
	"Fynance/internal/domain/login"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
)

type Handler struct {
	UserService  user.Service
	LoginService login.Service
	JwtService   utils.JwtService
}
