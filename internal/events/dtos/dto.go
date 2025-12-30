package dtos

import (
	"database/sql"
	"time"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusFinished  EventStatus = "finished"
)

type CreateRequest struct {
	Title         string      `json:"title" form:"title"`
	UserID        string      `json:"user_id" form:"user_id"`
	Description   string      `json:"description" form:"description"`
	StartDate     string      `json:"start_date" form:"start_date"`
	EndDate       string      `json:"end_date" form:"end_date"`
	City          string      `json:"city" form:"city"`
	Street        string      `json:"street" form:"street"`
	Neighbourhood string      `json:"neighbourhood" form:"neighbourhood"`
	MaxCapacity   int         `json:"max_capacity" form:"max_capacity"`
	Status        EventStatus `json:"status" form:"status"`
	PosterURL     string      `json:"poster_url" form:"poster_url"`
}

type UpdateRequest struct {
	ID            string      `json:"id" form:"id"`
	Title         string      `json:"title" form:"title"`
	UserID        string      `json:"user_id" form:"user_id"`
	Description   string      `json:"description,omitempty" form:"description"`
	StartDate     time.Time   `json:"start_date" form:"start_date"`
	EndDate       time.Time   `json:"end_date" form:"end_date"`
	City          string      `json:"city" form:"city"`
	Street        string      `json:"street" form:"street"`
	Neighbourhood string      `json:"neighbourhood" form:"neighbourhood"`
	MaxCapacity   int         `json:"max_capacity" form:"max_capacity"`
	Status        EventStatus `json:"status" form:"status"`
	PosterURL     string      `json:"poster_url,omitempty" form:"poster_url"`
}

type EventCard struct {
	ID          string
	UserID      string
	Title       string
	StartDate   string
	EndDate     string
	City        string
	MaxCapacity int64
	Status      string
	PosterURL   string
}

type Event struct {
	ID            string
	UserID        string
	Title         string
	Description   sql.NullString
	StartDate     time.Time
	EndDate       time.Time
	City          string
	Street        string
	Neighbourhood string
	MaxCapacity   int64
	Status        EventStatus
	PosterURL     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EventShow struct {
	ID            string
	UserID        string
	Title         string
	Description   *string
	StartDate     time.Time
	EndDate       time.Time
	City          string
	Street        string
	Neighbourhood string
	MaxCapacity   int64
	Status        EventStatus
	PosterURL     string
	Registrations int
	TicketBatch   TicketBatch
}

type TicketBatch struct {
	ID            string
	EventID       string
	BatchName     string
	Price         float64
	DateLimit     time.Time
	QuantityLimit int
	BatchOrder    int
	SoldQty       int
	TotalSold     float64
	Active        bool
}
