package contracts

import (
	"Fynance/internal/domain/transaction"
)

type TransactionCreateRequest struct {
	Type        string  `json:"type" binding:"required,oneof=RECEIPT EXPENSE TRANSFER GOALS INVESTMENT WITHDRAW"`
	CategoryID  string  `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=255"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon" binding:"omitempty,max=50"`
}

type TransactionCreateResponse struct {
	Message     string                  `json:"message"`
	Transaction transaction.Transaction `json:"transaction"`
}
