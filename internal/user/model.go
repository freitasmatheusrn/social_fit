package user

import "time"

type User struct {
	ID        string
	Name      string
	Cpf       string
	Email     string
	Phone     string
	BirthDate time.Time
	Password  string
	Admin     bool
	CreatedAt time.Time
}

type SignupRequest struct {
	Name                 string    `json:"name" form:"name"`
	Cpf                  string    `json:"cpf" form:"cpf"`
	Email                string    `json:"email" form:"email"`
	Phone                string    `json:"phone" form:"phone"`
	BirthDate            time.Time `json:"birth_date" form:"birth_date"`
	Password             string    `json:"password" form:"password"`
	PasswordConfirmation string    `json:"password_confirmation" form:"password_confirmation"`
}

type SignupResponse struct {
	ID    string `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Admin bool   `json:"admin" form:"admin"`
}

type SigninRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type SigninResponse struct {
	ID    string `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Admin bool   `json:"admin" form:"admin"`
}
