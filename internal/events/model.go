package events

import (
	"time"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusFinished  EventStatus = "finished"
)

type Event struct {
	ID            string
	UserID        string
	Title         string
	Description   string
	StartDate     time.Time
	EndDate       time.Time
	City          string
	Street        string
	Neighbourhood string
	MaxCapacity   int
	Status        EventStatus
	PosterURL     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
