package investment

import (
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"
	"crypto/rand"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository      Repository
	TransactionRepo transaction.Repository
}

func NewService(repo Repository, transactionRepo transaction.Repository) *Service {
	return &Service{
		Repository:      repo,
		TransactionRepo: transactionRepo,
	}
}

func (s *Service) CreateInvestment(req CreateInvestmentRequest) (*Investment, error) {
	investmentId := utils.GenerateULIDObject()
	investment := &Investment{
		Id:              investmentId,
		UserId:          req.UserId,
		Type:            req.Type,
		Name:            req.Name,
		CurrentBalance:  req.InitialAmount,
		ReturnRate:      req.ReturnRate,
		ApplicationDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.Repository.Create(investment); err != nil {
		return nil, err
	}

	transId := utils.GenerateULIDObject()

	trans := &transaction.Transaction{
		Id:           transId,
		UserId:       req.UserId,
		Type:         transaction.Investment,
		CategoryId:   req.CategoryId,
		Amount:       req.InitialAmount,
		Description:  "Aporte inicial - " + req.Name,
		Date:         time.Now(),
		InvestmentId: &investmentId,
	}

	if err := s.TransactionRepo.Create(trans); err != nil {
		s.Repository.Delete(investmentId, req.UserId)
		return nil, err
	}

	return investment, nil
}

func (s *Service) MakeContribution(investmentId, userId ulid.ULID, amount float64, categoryId ulid.ULID, description string) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	transId := utils.GenerateULIDObject()

	trans := &transaction.Transaction{
		Id:           transId,
		UserId:       userId,
		Type:         transaction.Investment,
		CategoryId:   categoryId,
		Amount:       amount,
		Description:  description,
		Date:         time.Now(),
		InvestmentId: &investmentId,
	}

	if err := s.TransactionRepo.Create(trans); err != nil {
		return err
	}

	investment.CurrentBalance += amount
	return s.Repository.Update(investment)
}

func (s *Service) MakeWithdraw(investmentId, userId ulid.ULID, amount float64, categoryId ulid.ULID, description string) error {
	investment, err := s.Repository.GetInvestmentById(investmentId, userId)
	if err != nil {
		return errors.New("investment not found")
	}

	if investment.CurrentBalance < amount {
		return errors.New("insufficient balance in investment")
	}

	entropy := ulid.Monotonic(rand.Reader, 0)
	transId := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	trans := &transaction.Transaction{
		Id:           transId,
		UserId:       userId,
		Type:         transaction.Withdraw,
		CategoryId:   categoryId,
		Amount:       amount,
		Description:  description,
		Date:         time.Now(),
		InvestmentId: &investmentId,
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
