package user

import (
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func (u User) ValidatePassword(password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), password)

}

func (u User) ValidateFields() error {
	err := EmailValid(u.Email)
	if err != nil {
		return err
	}
	err = PasswordValid(u.Password)
	if err != nil {
		return err
	}
	err = NameValid(u.Name)
	if err != nil {
		return err
	}
	err = CPFValid(u.Cpf)
	if err != nil {
		return err
	}
	err = BirthDateValid(u.BirthDate)
	if err != nil {
		return err
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

func PasswordValid(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("senha deve conter ao menos 8 caracteres")
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
		return fmt.Errorf("senha deve conter ao menos uma letra minúscula")
	}
	if !hasUpper {
		return fmt.Errorf("senha deve conter ao menos uma letra maiúscula")
	}
	if !hasNumber {
		return fmt.Errorf("senha deve conter ao menos um número")
	}
	if !hasSymbol {
		return fmt.Errorf("senha deve conter ao menos um símbolo")
	}

	return nil
}

func NameValid(name string) error {
	if len(name) <= 3 {
		return fmt.Errorf("nome deve conter ao menos 3 letras")
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
		return fmt.Errorf("cpf inválido")
	}

	allEqual := true
	for i := 1; i < 11; i++ {
		if digits[i] != digits[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return fmt.Errorf("cpf inválido")
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
		return fmt.Errorf("cpf inválido")
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
		return fmt.Errorf("cpf inválido")
	}

	return nil
}

func BirthDateValid(birthDate time.Time) error {
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
		return fmt.Errorf("idade mínima é 10 anos")
	}

	return nil
}
