package investment

import "github.com/oklog/ulid/v2"

type CreateInvestmentRequest struct {
	UserId        ulid.ULID `json:"user_id"`
	Type          Types     `json:"type"`
	Name          string    `json:"name"`
	InitialAmount float64   `json:"initial_amount"`
	ReturnRate    float64   `json:"return_rate"`
	CategoryId    ulid.ULID `json:"category_id"`
}

type ContributionRequest struct {
	Amount      float64   `json:"amount"`
	CategoryId  ulid.ULID `json:"category_id"`
	Description string    `json:"description"`
}

type WithdralRequest struct {
	Amount      float64   `json:"amount"`
	CategoryId  ulid.ULID `json:"category_id"`
	Description string    `json:"description"`
}
