package validation

import (
	"Fynance/internal/errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		Init()
	}

	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				return errors.NewValidationError(fieldError.Field(), getValidationMessage(fieldError))
			}
		}
		return errors.ErrValidation.WithError(err)
	}

	return nil
}

func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.NewValidationError("amount", "deve ser maior que zero")
	}
	return nil
}

func ValidateDate(date time.Time) error {
	if date.IsZero() {
		return errors.NewValidationError("date", "data é obrigatória")
	}
	return nil
}

func ValidateDateNotFuture(date time.Time) error {
	if err := ValidateDate(date); err != nil {
		return err
	}
	if date.After(time.Now()) {
		return errors.NewValidationError("date", "data não pode ser futura")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.NewValidationError("email", "email é obrigatório")
	}
	if err := validate.Var(email, "required,email"); err != nil {
		return errors.NewValidationError("email", "email inválido")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.NewValidationError("password", "senha deve ter pelo menos 8 caracteres")
	}
	return nil
}

func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "é obrigatório"
	case "email":
		return "deve ser um email válido"
	case "min":
		return "deve ter no mínimo " + fieldError.Param() + " caracteres"
	case "max":
		return "deve ter no máximo " + fieldError.Param() + " caracteres"
	case "gte":
		return "deve ser maior ou igual a " + fieldError.Param()
	case "lte":
		return "deve ser menor ou igual a " + fieldError.Param()
	case "gt":
		return "deve ser maior que " + fieldError.Param()
	case "lt":
		return "deve ser menor que " + fieldError.Param()
	default:
		return "inválido"
	}
}
