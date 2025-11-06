package investment

import (
	"context"

	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository      Repository
	TransactionRepo transaction.Repository
	UserService     *user.Service
}

func NewService(repo Repository, transactionRepo transaction.Repository) *Service {
	return &Service{
		Repository:      repo,
		TransactionRepo: transactionRepo,
	}
}

func (s *Service) CreateInvestment(ctx context.Context, req CreateInvestmentRequest) (*Investment, error) {
	investmentId := utils.GenerateULIDObject()

	if err := s.ensureUserExists(ctx, req.UserId); err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	investment, err := s.CreateInvestmentStruct(req, investmentId)
	if err != nil {
		return nil, err
	}

	if err = s.Repository.Create(ctx, investment); err != nil {
		return nil, err
	}

	trans, err := s.CreateTransactionStruct(req, investmentId)
	if err != nil {
		return nil, err
	}

	if err := s.TransactionRepo.Create(ctx, trans); err != nil {
		s.Repository.Delete(ctx, investmentId, req.UserId)
		return nil, err
	}

	return investment, nil
}

func (s *Service) MakeContribution(ctx context.Context, investmentId, userId ulid.ULID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	investment, err := s.Repository.GetInvestmentById(ctx, investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	trans, err := s.makeInvestmentMovement(investmentId, userId, amount, description, transaction.Investment)
	if err != nil {
		return err
	}

	if err := s.TransactionRepo.Create(ctx, trans); err != nil {
		return err
	}

	investment.CurrentBalance += amount
	return s.Repository.Update(ctx, investment)
}

func (s *Service) MakeWithdraw(ctx context.Context, investmentId, userId ulid.ULID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	investment, err := s.Repository.GetInvestmentById(ctx, investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	if investment.CurrentBalance < amount {
		return errors.New("insufficient balance in investment")
	}

	trans, err := s.makeInvestmentMovement(investmentId, userId, amount, description, transaction.Withdraw)
	if err != nil {
		return err
	}

	if err := s.TransactionRepo.Create(ctx, trans); err != nil {
		return err
	}

	investment.CurrentBalance -= amount
	return s.Repository.Update(ctx, investment)
}

func (s *Service) ListInvestments(ctx context.Context, userId ulid.ULID) ([]*Investment, error) {
	if err := s.ensureUserExists(ctx, userId); err != nil {
		return nil, err
	}
	return s.Repository.GetByUserId(ctx, userId)
}

func (s *Service) GetInvestment(ctx context.Context, investmentId, userId ulid.ULID) (*Investment, error) {
	if err := s.ensureUserExists(ctx, userId); err != nil {
		return nil, err
	}
	return s.Repository.GetInvestmentById(ctx, investmentId, userId)
}

func (s *Service) GetTotalInvested(ctx context.Context, investmentId, userId ulid.ULID) (float64, error) {
	transactions, err := s.TransactionRepo.GetByInvestmentId(ctx, investmentId, userId)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, trans := range transactions {
		switch trans.Type {
		case transaction.Investment:
			total += trans.Amount
		case transaction.Withdraw:
			total -= trans.Amount
		}
	}

	return total, nil
}

func (s *Service) CalculateReturn(ctx context.Context, investmentId, userId ulid.ULID) (float64, float64, error) {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentId, userId)
	if err != nil {
		return 0, 0, err
	}

	totalInvested, err := s.GetTotalInvested(ctx, investmentId, userId)
	if err != nil {
		return 0, 0, err
	}

	if totalInvested == 0 {
		return 0, 0, nil
	}

	profit := investment.CurrentBalance - totalInvested
	returnPercentage := (profit / totalInvested) * 100

	return profit, returnPercentage, nil
}

func (s *Service) DeleteInvestment(ctx context.Context, investmentId, userId ulid.ULID) error {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentId, userId)
	if err != nil {
		return err
	}

	if investment.CurrentBalance > 0 {
		return errors.New("cannot delete investment with balance")
	}

	return s.Repository.Delete(ctx, investmentId, userId)
}

func (s *Service) UpdateInvestment(ctx context.Context, investmentId, userId ulid.ULID, req UpdateInvestmentRequest) error {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentId, userId)
	if err != nil {
		return err
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return errors.New("name is required")
		}
		investment.Name = trimmed
	}

	if req.Type != nil && *req.Type != "" {
		investment.Type = Types(*req.Type)
	}

	if req.ReturnRate != nil {
		investment.ReturnRate = *req.ReturnRate
	}

	investment.UpdatedAt = time.Now()

	return s.Repository.Update(ctx, investment)
}

func (s *Service) CreateInvestmentStruct(req CreateInvestmentRequest, InvestmentId ulid.ULID) (*Investment, error) {
	now := utils.SetTimestamps()

	investment := &Investment{
		Id:              InvestmentId,
		UserId:          req.UserId,
		Type:            Types(req.Type),
		Name:            req.Name,
		CurrentBalance:  req.InitialAmount,
		ReturnBalance:   0,
		ReturnRate:      req.ReturnRate,
		ApplicationDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	return investment, nil
}

func (s *Service) CreateTransactionStruct(req CreateInvestmentRequest, InvestmentId ulid.ULID) (*transaction.Transaction, error) {
	now := utils.SetTimestamps()

	return &transaction.Transaction{
		Id:           utils.GenerateULIDObject(),
		UserId:       req.UserId,
		Type:         transaction.Investment,
		Amount:       req.InitialAmount,
		Description:  "Aporte inicial - " + req.Name,
		Date:         now,
		InvestmentId: &InvestmentId,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *Service) makeInvestmentMovement(investmentId, userId ulid.ULID, amount float64, description string, movementType transaction.Types) (*transaction.Transaction, error) {
	desc := strings.TrimSpace(description)
	if desc == "" {
		if movementType == transaction.Withdraw {
			desc = "Resgate"
		} else {
			desc = "Aporte"
		}
	}

	now := utils.SetTimestamps()

	return &transaction.Transaction{
		Id:           utils.GenerateULIDObject(),
		UserId:       userId,
		Type:         movementType,
		Amount:       amount,
		Description:  desc,
		Date:         now,
		InvestmentId: &investmentId,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (s *Service) ensureUserExists(ctx context.Context, userID ulid.ULID) error {
	if s.UserService == nil {
		return errors.New("user service not configured")
	}
	_, err := s.UserService.GetByID(ctx, userID.String())
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}
