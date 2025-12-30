package events

import (
	"context"
	"time"

	"github.com/freitasmatheusrn/social-fit/internal/events/dtos"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
)

type Service struct {
	repo RepositoryInterface
}

type ServiceInterface interface {
	Create(ctx context.Context, event dtos.CreateRequest) (string, *rest.ApiErr)
	Delete(ctx context.Context, eventID string, userID string) *rest.ApiErr
	Update(ctx context.Context, event dtos.UpdateRequest) *rest.ApiErr
	FindManyByUser(ctx context.Context, userID string) ([]dtos.Event, *rest.ApiErr)
	GetEventWithDetails(ctx context.Context, eventID, userID string) (*dtos.EventShow, error) 
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, event dtos.CreateRequest) (string, *rest.ApiErr) {
	startDate, err := time.Parse("2006-01-02T15:04", event.StartDate)
	if err != nil {
		return "", rest.NewBadRequestError("data de ínicio inválida")
	}

	endDate, err := time.Parse("2006-01-02T15:04", event.EndDate)
	if err != nil {
		return "", rest.NewBadRequestError("data de término inválida")
	}
	e := Event{
		UserID:        event.UserID,
		PosterURL:     event.PosterURL,
		Title:         event.Title,
		Description:   event.Description,
		StartDate:     startDate,
		EndDate:       endDate,
		City:          event.City,
		Street:        event.Street,
		Neighbourhood: event.Neighbourhood,
		MaxCapacity:   event.MaxCapacity,
		Status:        EventStatus(event.Status),
	}
	if appErr := e.ValidateFields(); appErr != nil {
		return "", appErr
	}
	id, appErr := s.repo.Create(ctx, e)
	if appErr != nil {
		return "", appErr
	}
	return id, nil
}

func (s *Service) Delete(ctx context.Context, eventID string, userID string) *rest.ApiErr {
	if appErr := s.repo.Delete(ctx, eventID, userID); appErr != nil {
		return appErr
	}
	return nil
}

func (s *Service) Update(ctx context.Context, event dtos.UpdateRequest) *rest.ApiErr {
	e := Event{
		ID:            event.ID,
		UserID:        event.UserID,
		PosterURL:     event.PosterURL,
		Title:         event.Title,
		Description:   event.Description,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		City:          event.City,
		Street:        event.Street,
		Neighbourhood: event.Neighbourhood,
		MaxCapacity:   event.MaxCapacity,
		Status:        EventStatus(event.Status),
	}
	if err := e.ValidateFields(); err != nil {
		return rest.NewBadRequestValidationError("Campo(s) inválidos", err.Causes)
	}
	if err := s.repo.Update(ctx, e); err != nil {
		return err
	}
	return nil
}

func (s *Service) FindManyByUser(ctx context.Context, userID string) ([]dtos.Event, *rest.ApiErr) {
	events, err := s.repo.FindManyByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) GetEventWithDetails(ctx context.Context, eventID, userID string) (*dtos.EventShow, error) {
	event, err := s.repo.GetEventWithDetails(ctx, eventID, userID)
	if err != nil{
		return nil, err
	}
	return event, nil
}