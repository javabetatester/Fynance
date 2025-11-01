package routes

import (
	login "Fynance/internal/domain/auth"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
)

type Handler struct {
	UserService        user.Service
	LoginService       login.Service
	JwtService         *utils.JwtService
	TransactionService transaction.Service
}
