package user

import (
	"context"

	"github.com/freitasmatheusrn/social-fit/internal/events"
	evtDtos "github.com/freitasmatheusrn/social-fit/internal/events/dtos"
	"github.com/freitasmatheusrn/social-fit/pkg/fmtdate"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
)

type Service struct {
	repo         *Repository
	eventService events.ServiceInterface
}
type ServiceInterface interface {
	Login(ctx context.Context, credentials SigninRequest) (SigninResponse, *rest.ApiErr)
	Signup(ctx context.Context, user SignupRequest) (SignupResponse, *rest.ApiErr)
	Home(ctx context.Context, userID string) (evtDtos.EventCard, *rest.ApiErr)
}

func NewService(repo *Repository, e events.ServiceInterface) *Service {
	return &Service{
		repo:         repo,
		eventService: e,
	}
}

func (s *Service) Login(ctx context.Context, credentials SigninRequest) (SigninResponse, *rest.ApiErr) {
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

func (s *Service) Signup(ctx context.Context, request SignupRequest) (SignupResponse, *rest.ApiErr) {
	u := User{
		Name:      request.Name,
		Email:     request.Email,
		Cpf:       request.Cpf,
		Phone:     request.Phone,
		BirthDate: request.BirthDate,
		Password:  request.Password,
	}
	err := u.ValidateFields()
	if err != nil {
		return SignupResponse{}, err
	}
	user, err := s.repo.Signup(ctx, u)
	if err != nil {
		return SignupResponse{}, err
	}
	return SignupResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Admin: user.Admin,
	}, nil
}

func (s *Service) Home(ctx context.Context, userID string) ([]evtDtos.EventCard, *rest.ApiErr) {
	events, err := s.eventService.FindManyByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var e []evtDtos.EventCard
	for _, event := range events {
		e = append(e, evtDtos.EventCard{
			ID:          event.ID,
			UserID:      event.UserID,
			Title:       event.Title,
			StartDate:   event.StartDate.Format(fmtdate.DateTimeLayout),
			EndDate:     event.EndDate.Format(fmtdate.DateTimeLayout),
			City:        event.City,
			MaxCapacity: event.MaxCapacity,
			Status:      string(event.Status),
		})
	}
	return e, nil
}
