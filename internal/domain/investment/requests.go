package investment

import "github.com/oklog/ulid/v2"

type CreateInvestmentRequest struct {
	UserId        ulid.ULID `json:"user_id"`
	Type          Types     `json:"type"`
	Name          string    `json:"name"`
	InitialAmount float64   `json:"initial_amount"`
	ReturnRate    float64   `json:"return_rate"`
}

type ContributionRequest struct {
	UserId      ulid.ULID `json:"user_id"`
	Id          ulid.ULID `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type WithdralRequest struct {
	UserId      ulid.ULID `json:"user_id"`
	Id          ulid.ULID `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type UpdateInvestmentRequest struct {
	UserId        ulid.ULID `json:"user_id"`
	Id            ulid.ULID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type" binding:"omitempty,oneof=CDB LCI LCA TESOURO_DIRETO ACOES FUNDOS CRIPTOMOEDAS PREVIDENCIA"`
	InitialAmount float64   `json:"initial_amount" binding:"omitempty"`
	ReturnRate    float64   `json:"return_rate" binding:"omitempty"`
}
