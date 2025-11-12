package investment

import (
	"context"
	"fmt"
	"strings"
	"time"

	domaincontracts "Fynance/internal/domain/contracts"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository      Repository
	TransactionRepo transaction.Repository
	UserService     *user.Service
}

func NewService(repo Repository, transactionRepo transaction.Repository) *Service {
	return &Service{Repository: repo, TransactionRepo: transactionRepo}
}

func (s *Service) CreateInvestment(ctx context.Context, req domaincontracts.CreateInvestmentRequest) (*Investment, error) {
	if err := s.ensureUserExists(ctx, req.UserId); err != nil {
		return nil, err
	}

	trimmedName := strings.TrimSpace(req.Name)
	if trimmedName == "" {
		return nil, appErrors.NewValidationError("name", "é obrigatório")
	}
	req.Name = trimmedName

	investmentID := pkg.GenerateULIDObject()
	entity := s.CreateInvestmentStruct(req, investmentID)

	if err := s.Repository.Create(ctx, entity); err != nil {
		return nil, err
	}

	movement := s.CreateTransactionStruct(req, investmentID)
	if err := s.TransactionRepo.Create(ctx, movement); err != nil {
		_ = s.Repository.Delete(ctx, investmentID, req.UserId)
		return nil, err
	}

	return entity, nil
}

func (s *Service) MakeContribution(ctx context.Context, investmentID, userID ulid.ULID, amount float64, description string) error {
	if amount <= 0 {
		return appErrors.NewValidationError("amount", "deve ser maior que zero")
	}

	investment, err := s.Repository.GetInvestmentById(ctx, investmentID, userID)
	if err != nil {
		return err
	}

	movement := s.makeInvestmentMovement(investmentID, userID, amount, description, transaction.Investment)
	if err := s.TransactionRepo.Create(ctx, movement); err != nil {
		return err
	}

	investment.CurrentBalance += amount
	return s.Repository.Update(ctx, investment)
}

func (s *Service) MakeWithdraw(ctx context.Context, investmentID, userID ulid.ULID, amount float64, description string) error {
	if amount <= 0 {
		return appErrors.NewValidationError("amount", "deve ser maior que zero")
	}

	investment, err := s.Repository.GetInvestmentById(ctx, investmentID, userID)
	if err != nil {
		return err
	}

	if investment.CurrentBalance < amount {
		return appErrors.NewValidationError("amount", "saldo insuficiente no investimento")
	}

	movement := s.makeInvestmentMovement(investmentID, userID, amount, description, transaction.Withdraw)
	if err := s.TransactionRepo.Create(ctx, movement); err != nil {
		return err
	}

	investment.CurrentBalance -= amount
	return s.Repository.Update(ctx, investment)
}

func (s *Service) ListInvestments(ctx context.Context, userID ulid.ULID) ([]*Investment, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}
	return s.Repository.GetByUserId(ctx, userID)
}

func (s *Service) GetInvestment(ctx context.Context, investmentID, userID ulid.ULID) (*Investment, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}
	return s.Repository.GetInvestmentById(ctx, investmentID, userID)
}

func (s *Service) GetTotalInvested(ctx context.Context, investmentID, userID ulid.ULID) (float64, error) {
	transactions, err := s.TransactionRepo.GetByInvestmentId(ctx, investmentID, userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, tx := range transactions {
		switch tx.Type {
		case transaction.Investment:
			total += tx.Amount
		case transaction.Withdraw:
			total -= tx.Amount
		}
	}

	return total, nil
}

func (s *Service) CalculateReturn(ctx context.Context, investmentID, userID ulid.ULID) (float64, float64, error) {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentID, userID)
	if err != nil {
		return 0, 0, err
	}

	totalInvested, err := s.GetTotalInvested(ctx, investmentID, userID)
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

func (s *Service) DeleteInvestment(ctx context.Context, investmentID, userID ulid.ULID) error {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentID, userID)
	if err != nil {
		return err
	}

	if investment.CurrentBalance > 0 {
		return appErrors.NewValidationError("investment", "possui saldo e não pode ser removido")
	}

	return s.Repository.Delete(ctx, investmentID, userID)
}

func (s *Service) UpdateInvestment(ctx context.Context, investmentID, userID ulid.ULID, req domaincontracts.UpdateInvestmentRequest) error {
	investment, err := s.Repository.GetInvestmentById(ctx, investmentID, userID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return appErrors.NewValidationError("name", "é obrigatório")
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

func (s *Service) CreateInvestmentStruct(req domaincontracts.CreateInvestmentRequest, investmentID ulid.ULID) *Investment {
	now := pkg.SetTimestamps()

	return &Investment{
		Id:              investmentID,
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
}

func (s *Service) CreateTransactionStruct(req domaincontracts.CreateInvestmentRequest, investmentID ulid.ULID) *transaction.Transaction {
	now := pkg.SetTimestamps()

	return &transaction.Transaction{
		Id:           pkg.GenerateULIDObject(),
		UserId:       req.UserId,
		Type:         transaction.Investment,
		Amount:       req.InitialAmount,
		Description:  "Aporte inicial - " + req.Name,
		Date:         now,
		InvestmentId: &investmentID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *Service) makeInvestmentMovement(investmentID, userID ulid.ULID, amount float64, description string, movementType transaction.Types) *transaction.Transaction {
	desc := strings.TrimSpace(description)
	if desc == "" {
		if movementType == transaction.Withdraw {
			desc = "Resgate"
		} else {
			desc = "Aporte"
		}
	}

	now := pkg.SetTimestamps()

	return &transaction.Transaction{
		Id:           pkg.GenerateULIDObject(),
		UserId:       userID,
		Type:         movementType,
		Amount:       amount,
		Description:  desc,
		Date:         now,
		InvestmentId: &investmentID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *Service) ensureUserExists(ctx context.Context, userID ulid.ULID) error {
	if s.UserService == nil {
		return appErrors.ErrInternalServer.WithError(fmt.Errorf("serviço de usuário não configurado"))
	}
	_, err := s.UserService.GetByID(ctx, userID.String())
	if err != nil {
		return appErrors.ErrUserNotFound.WithError(err)
	}
	return nil
}
