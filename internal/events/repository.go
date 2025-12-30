package events

import (
	"context"
	"time"

	"github.com/freitasmatheusrn/social-fit/internal/events/dtos"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

type RepositoryInterface interface {
	Create(ctx context.Context, event Event) (string, *rest.ApiErr)
	Delete(ctx context.Context, eventID string, userID string) *rest.ApiErr
	Update(ctx context.Context, event Event) *rest.ApiErr
	FindManyByUser(ctx context.Context, userID string) ([]dtos.Event, *rest.ApiErr)
	GetEventWithDetails(ctx context.Context, eventID, userID string) (*dtos.EventShow, error) 
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, event Event) (string, *rest.ApiErr) {
	query := `
		INSERT INTO events (
			title, description, start_date, end_date, 
			city, street, neighbourhood, max_capacity, 
			status, user_id, created_at, updated_at, poster_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), $11)
		RETURNING id
	`

	var id string
	err := r.db.QueryRow(
		ctx,
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.City,
		event.Street,
		event.Neighbourhood,
		event.MaxCapacity,
		event.Status,
		event.UserID,
		event.PosterURL,
	).Scan(&id)

	if err != nil {
		return "", rest.NewInternalServerError("falha ao criar o evento " + err.Error())
	}

	return id, nil
}

func (r *Repository) Delete(ctx context.Context, eventID string, userID string) *rest.ApiErr {
	query := `
		DELETE FROM events 
		WHERE id = $1 AND user_id = $2
	`

	cmdTag, err := r.db.Exec(ctx, query, eventID, userID)
	if err != nil {
		return rest.NewInternalServerError("falha ao deletar o evento")
	}

	if cmdTag.RowsAffected() == 0 {
		return rest.NewForbiddenError("evento não encontrado ou não autorizado")
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, event Event) *rest.ApiErr {
	query := `
		UPDATE events 
		SET 
			title = $1,
			description = $2,
			start_date = $3,
			end_date = $4,
			city = $5,
			street = $6,
			neighbourhood = $7,
			max_capacity = $8,
			status = $9,
			updated_at = NOW()
			poster_url = $10
		WHERE id = $11 AND user_id = $12
	`

	cmdTag, err := r.db.Exec(
		ctx,
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.City,
		event.Street,
		event.Neighbourhood,
		event.MaxCapacity,
		event.Status,
		event.ID,
		event.UserID,
	)

	if err != nil {
		return rest.NewInternalServerError("falha ao atualizar o evento")
	}

	if cmdTag.RowsAffected() == 0 {
		return rest.NewForbiddenError("evento não encontrado ou não autorizado")
	}

	return nil
}

func (r *Repository) FindManyByUser(ctx context.Context, userID string) ([]dtos.Event, *rest.ApiErr) {
	query := `
		SELECT 
			id, user_id, title, description, start_date, end_date,
			city, street, neighbourhood, max_capacity, status,
			created_at, updated_at
		FROM events
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, rest.NewInternalServerError("falha ao buscar eventos")
	}
	defer rows.Close()

	var events []dtos.Event
	for rows.Next() {
		var event dtos.Event
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.City,
			&event.Street,
			&event.Neighbourhood,
			&event.MaxCapacity,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, rest.NewInternalServerError("falha ao parsear dados")
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, rest.NewInternalServerError("falha na iteração dos dados")
	}

	return events, nil
}

func (r *Repository) GetEventWithDetails(ctx context.Context, eventID, userID string) (*dtos.EventShow, error) {
    query := `
        WITH event_data AS (
            SELECT 
                e.id,
                e.user_id,
                e.title,
                e.description,
                e.start_date,
                e.end_date,
                e.city,
                e.street,
                e.neighbourhood,
                e.max_capacity,
                e.status,
                e.poster_url,
                COUNT(DISTINCT r.id) FILTER (WHERE r.deleted_at IS NULL) as total_registrations
            FROM events e
            LEFT JOIN registrations r ON e.id = r.event_id
            WHERE e.id = $1 AND e.user_id = $2
            GROUP BY e.id
        ),
        batch_data AS (
            SELECT 
                tb.id,
                tb.event_id,
                tb.batch_name,
                tb.price,
                tb.date_limit,
                tb.quantity_limit,
                tb.batch_order,
                tb.active,
                COUNT(r.id) FILTER (WHERE r.deleted_at IS NULL AND r.payment_status = 'confirmed') as sold_qty,
                COALESCE(SUM(r.amount_paid) FILTER (WHERE r.deleted_at IS NULL AND r.payment_status = 'confirmed'), 0) as total_sold
            FROM ticket_batches tb
            LEFT JOIN registrations r ON tb.id = r.ticket_batch_id
            WHERE tb.event_id = $1 
                AND tb.deleted_at IS NULL
                AND tb.active = true
            GROUP BY tb.id
        )
        SELECT 
            ed.id,
            ed.user_id,
            ed.title,
            ed.description,
            ed.start_date,
            ed.end_date,
            ed.city,
            ed.street,
            ed.neighbourhood,
            ed.max_capacity,
            ed.status,
            ed.poster_url,
            ed.total_registrations,
            bd.id as batch_id,
            bd.event_id as batch_event_id,
            bd.batch_name,
            bd.price,
            bd.date_limit,
            bd.quantity_limit,
            bd.batch_order,
            bd.sold_qty,
            bd.total_sold,
            bd.active
        FROM event_data ed
        LEFT JOIN batch_data bd ON ed.id = bd.event_id
    `

    var event dtos.EventShow
    var batch dtos.TicketBatch
    
    var batchID, batchEventID, batchName *string
    var price, totalSold *float64
    var dateLimit *time.Time
    var quantityLimit, batchOrder, soldQty *int
    var active *bool

    err := r.db.QueryRow(ctx, query, eventID, userID).Scan(
        &event.ID,
        &event.UserID,
        &event.Title,
        &event.Description,
        &event.StartDate,
        &event.EndDate,
        &event.City,
        &event.Street,
        &event.Neighbourhood,
        &event.MaxCapacity,
        &event.Status,
        &event.PosterURL,
        &event.Registrations,
        &batchID,
        &batchEventID,
        &batchName,
        &price,
        &dateLimit,
        &quantityLimit,
        &batchOrder,
        &soldQty,
        &totalSold,
        &active,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, rest.NewNotFoundError("evento não encontrado")
        }
        return nil, rest.NewInternalServerError("erro interno do servidor")
    }

    // Preenche o ticket batch se existir
    if batchID != nil {
        batch.ID = *batchID
        batch.EventID = *batchEventID
        batch.BatchName = *batchName
        batch.Price = *price
        batch.DateLimit = *dateLimit
        batch.QuantityLimit = *quantityLimit
        batch.BatchOrder = *batchOrder
        batch.SoldQty = *soldQty
        batch.TotalSold = *totalSold
        batch.Active = *active
        
        event.TicketBatch = batch
    }

    return &event, nil
}
