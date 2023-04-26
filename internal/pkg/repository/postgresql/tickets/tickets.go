package tickets

import (
	"context"
	"database/sql"
	"time"

	"homework-week-5/internal/pkg/db"
	"homework-week-5/internal/pkg/repository"
)

type TicketRepo struct {
	db db.DBops
}

func NewTickets(db db.DBops) *TicketRepo {
	return &TicketRepo{db: db}
}

func (r *TicketRepo) Add(ctx context.Context, ticket *repository.Ticket) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO tickets(user_id, cost, place) VALUES ($1, $2, $3) RETURNING id`,
		ticket.UserID, ticket.Cost, ticket.Place).Scan(&id)
	return id, err
}

func (r *TicketRepo) GetById(ctx context.Context, id int) (*repository.Ticket, error) {
	var u repository.Ticket
	err := r.db.Get(ctx, &u, "SELECT id, user_id, cost, place, created_at, updated_at FROM tickets WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &u, err
}

func (r *TicketRepo) GetByUserId(ctx context.Context, UserId int) ([]*repository.Ticket, error) {
	tickets := make([]*repository.Ticket, 0)
	err := r.db.Select(ctx, &tickets, "SELECT id, user_id, cost, place, created_at, updated_at FROM tickets WHERE user_id=$1", UserId)
	return tickets, err
}

func (r *TicketRepo) List(ctx context.Context) ([]*repository.Ticket, error) {
	tickets := make([]*repository.Ticket, 0)
	err := r.db.Select(ctx, &tickets, "SELECT id, user_id, cost, place, created_at, updated_at FROM tickets")
	return tickets, err
}

func (r *TicketRepo) Update(ctx context.Context, ticket *repository.Ticket) (bool, error) {
	start := time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE tickets SET user_id = $1, cost = $2, place = $3, updated_at = $4 WHERE id = $5",
		ticket.UserID, ticket.Cost, ticket.Place, start, ticket.ID)
	return result.RowsAffected() > 0, err
}

func (r *TicketRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM tickets WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
