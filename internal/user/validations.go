package user

import (
	"errors"
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"golang.org/x/crypto/bcrypt"
)

func (u User) ComparePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), password)

}

func (u User) ValidateFields() *rest.ApiErr {
	var causes []rest.Causes

	if err := EmailValid(u.Email); err != nil {
		causes = append(causes, rest.Causes{Field: "email", Message: err.Error()})
	}
	if passwordErrs := PasswordValid(u.Password); passwordErrs != nil {
		for _, e := range passwordErrs {
			causes = append(causes, rest.Causes{Field: "password", Message: e.Error()})
		}
	}
	if err := NameValid(u.Name); err != nil {
		causes = append(causes, rest.Causes{Field: "name", Message: err.Error()})
	}
	if err := CPFValid(u.Cpf); err != nil {
		causes = append(causes, rest.Causes{Field: "cpf", Message: err.Error()})
	}
	if birthDateErrs := BirthDateValid(u.BirthDate); birthDateErrs != nil {
		for _, e := range birthDateErrs {
			causes = append(causes, rest.Causes{Field: "birth_date", Message: e.Error()})
		}
	}
	if causes != nil {
		return rest.NewBadRequestValidationError("Campo(s) inválidos", causes)
	}
	return nil
}

func EmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("inválido %s", err)
	}
	return nil
}

func PasswordValid(password string) []error {
	var errs []error
	if len(password) < 8 {
		errs = append(errs, errors.New("senha precisa ter 8 ou mais caracteres"))
	}
	var hasLower, hasUpper, hasNumber, hasSymbol bool

	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}
	if !hasLower {
		errs = append(errs, errors.New("senha deve conter ao menos uma letra minúscula"))
	}
	if !hasUpper {
		errs = append(errs, errors.New("senha deve conter ao menos uma letra maiúscula"))

	}
	if !hasNumber {
		errs = append(errs, errors.New("senha deve conter ao menos um número"))
	}
	if !hasSymbol {
		errs = append(errs, errors.New("senha deve conter ao menos um símbolo"))
	}
	if errs != nil {
		return errs
	}
	return nil
}

func NameValid(name string) error {
	if len(name) <= 3 {
		return errors.New("nome deve conter ao menos 3 letras")
	}
	return nil
}

func CPFValid(cpf string) error {
	var digits []int
	for _, r := range cpf {
		if unicode.IsDigit(r) {
			digits = append(digits, int(r-'0'))
		}
	}

	if len(digits) != 11 {
		return errors.New("cpf inválido")
	}

	allEqual := true
	for i := 1; i < 11; i++ {
		if digits[i] != digits[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return errors.New("cpf inválido")
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	firstDV := (sum * 10) % 11
	if firstDV == 10 {
		firstDV = 0
	}
	if firstDV != digits[9] {
		return errors.New("cpf inválido")
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	secondDV := (sum * 10) % 11
	if secondDV == 10 {
		secondDV = 0
	}
	if secondDV != digits[10] {
		return errors.New("cpf inválido")
	}

	return nil
}

func BirthDateValid(birthDateStr string) []error {
	var errs []error
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		errs = append(errs, err)
	}
	now := time.Now()

	birth := time.Date(
		birthDate.Year(),
		birthDate.Month(),
		birthDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	age := today.Year() - birth.Year()

	if today.Month() < birth.Month() ||
		(today.Month() == birth.Month() && today.Day() < birth.Day()) {
		age--
	}

	if age < 10 {
		errs = append(errs, errors.New("idade mínima é 10 anos"))
	}
	if errs != nil {
		return errs
	}
	return nil
}
