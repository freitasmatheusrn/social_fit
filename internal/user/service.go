package user

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}
type ServiceInterface interface {
	Login(ctx context.Context, credentials SigninRequest) (SigninResponse, error)
	Signup(ctx context.Context, user SignupRequest) (SignupResponse, error)
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Login(ctx context.Context, credentials SigninRequest) (SigninResponse, error) {
	user := User{
		Email:    credentials.Email,
		Password: credentials.Password,
	}
	u, err := s.repo.Login(ctx, user)
	if err != nil {
		return SigninResponse{}, err
	}
	return SigninResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Admin: u.Admin,
	}, nil

}

func (s *Service) Signup(ctx context.Context, request SignupRequest) (SignupResponse, error) {
	if request.Password != request.PasswordConfirmation{
		return SignupResponse{}, fmt.Errorf("as senhas não estão iguais")
	}
	u := User{
		Name: request.Name,
		Email: request.Email,
		Cpf: request.Cpf,
		Phone: request.Phone,
		BirthDate: request.BirthDate,
		Password: request.Password,
	}
	err := u.ValidateFields()
	if err != nil {
		return SignupResponse{}, err
	}
	user, err := s.repo.Signup(ctx, u)
	if err != nil{
		return SignupResponse{}, err
	}
	return SignupResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Admin: user.Admin,
	}, nil
}
