package auth

import (
	"Fynance/internal/domain/user"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository  user.Repository
	UserService *user.Service
}

func (s *Service) Login(login Login) (*user.User, error) {
	if !s.UserExists(login.Email) {
		return nil, errors.New("account does not exist")
	}

	user, err := s.GetByEmail(login.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := PasswordValidate(login.Password, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Register(user *user.User) error {
	if s.UserExists(user.Email) {
		return errors.New("email already registered")
	}

	if err := PasswordRequirements(user.Password); err != nil {
		return err
	}

	// Use UserService.Create instead of Repository.Create to ensure ID generation
	if err := s.UserService.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByEmail(email string) (*user.User, error) {
	user, err := s.Repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UserExists(email string) bool {
	_, err := s.Repository.GetByEmail(email)
	return err == nil
}

func PasswordRequirements(senha string) error {
	if len(senha) < 8 {
		return errors.New("the password must be at least 8 characters long")
	}

	temMaiuscula, _ := regexp.MatchString(`[A-Z]`, senha)
	if !temMaiuscula {
		return errors.New("the password must contain at least one uppercase letter")
	}

	temEspecial, _ := regexp.MatchString(`[@$!%*?&#]`, senha)
	if !temEspecial {
		return errors.New("the password must contain at least one special character (@$!%*?&#)")
	}

	return nil
}

func PasswordValidate(lpassword string, upassword string) error {

	if lpassword == "" {
		return errors.New("password must be filled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(upassword), []byte(lpassword)); err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}

func PasswordHashing(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	password = string(hash)
	return password, nil
}
