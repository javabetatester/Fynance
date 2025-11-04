package investment

import (
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"errors"
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

func (s *Service) CreateInvestment(req CreateInvestmentRequest) (*Investment, error) {
	investmentId := utils.GenerateULIDObject()

	_, err := s.UserService.GetByID(req.UserId.String())
	if err != nil {
		return nil, err
	}

	investment, err := s.CreateInvestmentStruct(req, investmentId)
	if err != nil {
		return nil, err
	}

	if err = s.Repository.Create(investment); err != nil {
		return nil, err
	}

	trans, err := s.CreateTransactionStruct(req, investmentId)
	if err != nil {
		return nil, err
	}

	if err := s.TransactionRepo.Create(trans); err != nil {
		s.Repository.Delete(investmentId, req.UserId)
		return nil, err
	}

	return investment, nil
}

func (s *Service) MakeContribution(investmentId, userId ulid.ULID, amount float64, description string) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	trans, err := s.MakeWithdrawAndContributionStruct(investmentId, userId, amount, description)
	if err != nil {
		return err
	}

	if err := s.TransactionRepo.Create(trans); err != nil {
		return err
	}

	investment.CurrentBalance += amount
	return s.Repository.Update(investment)
}

func (s *Service) MakeWithdraw(investmentId, userId ulid.ULID, amount float64, description string) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	if investment.CurrentBalance < amount {
		return errors.New("insufficient balance in investment")
	}

	trans, err := s.MakeWithdrawAndContributionStruct(investmentId, userId, amount, description)
	if err != nil {
		return err
	}

	if err := s.TransactionRepo.Create(trans); err != nil {
		return err
	}

	investment.CurrentBalance -= amount
	return s.Repository.Update(investment)
}

func (s *Service) ListInvestments(userId ulid.ULID) ([]*Investment, error) {
	return s.Repository.GetByUserId(userId)
}

func (s *Service) GetInvestment(investmentId, userId ulid.ULID) (*Investment, error) {
	return s.Repository.GetInvestmentById(investmentId, userId)
}

func (s *Service) GetTotalInvested(investmentId, userId ulid.ULID) (float64, error) {
	transactions, err := s.TransactionRepo.GetByInvestmentId(investmentId, userId)
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

func (s *Service) CalculateReturn(investmentId, userId ulid.ULID) (float64, float64, error) {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return 0, 0, err
	}

	totalInvested, err := s.GetTotalInvested(investmentId, userId)
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

func (s *Service) DeleteInvestment(investmentId, userId ulid.ULID) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return err
	}

	if investment.CurrentBalance > 0 {
		return errors.New("cannot delete investment with balance")
	}

	return s.Repository.Delete(investmentId, userId)
}

func (s *Service) UpdateInvestment(investmentId, userId ulid.ULID, req UpdateInvestmentRequest) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return err
	}

	investment.Name = req.Name
	investment.Type = Types(req.Type)
	investment.ReturnRate = req.ReturnRate

	return s.Repository.Update(investment)
}

func (s *Service) CreateInvestmentStruct(req CreateInvestmentRequest, InvestmentId ulid.ULID) (*Investment, error) {
	investment := &Investment{
		Id:             InvestmentId,
		UserId:         req.UserId,
		Type:           Types(req.Type),
		Name:           req.Name,
		CurrentBalance: req.InitialAmount,
		ReturnRate:     req.ReturnRate,
	}
	now := utils.SetTimestamps()
	investment.CreatedAt = now
	investment.UpdatedAt = now
	return investment, nil
}

func (s *Service) CreateTransactionStruct(req CreateInvestmentRequest, InvestmentId ulid.ULID) (*transaction.Transaction, error) {

	transaction := &transaction.Transaction{
		Id:           utils.GenerateULIDObject(),
		UserId:       req.UserId,
		Type:         transaction.Investment,
		Amount:       req.InitialAmount,
		Description:  "Aporte inicial - " + req.Name,
		Date:         time.Now(),
		InvestmentId: &InvestmentId,
	}
	now := utils.SetTimestamps()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
	return transaction, nil
}

func (s *Service) MakeWithdrawAndContributionStruct(investmentId, userId ulid.ULID, amount float64, description string) (*transaction.Transaction, error) {
	transaction := &transaction.Transaction{
		Id:           utils.GenerateULIDObject(),
		UserId:       userId,
		Type:         transaction.Investment,
		Amount:       amount,
		Description:  description,
		Date:         time.Now(),
		InvestmentId: &investmentId,
	}
	now := utils.SetTimestamps()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
	return transaction, nil
}
