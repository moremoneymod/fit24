package repository

import (
	"context"

	"fit24/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) domain.LeadRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) SaveOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (created_at, plan, name, phone, email)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		order.CreatedAt,
		order.Plan,
		order.Name,
		order.Phone,
		order.Email,
	)
	return err
}

func (r *postgresRepository) SaveContact(ctx context.Context, contact *domain.Contact) error {
	query := `
		INSERT INTO contacts (created_at, name, phone, message)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query,
		contact.CreatedAt,
		contact.Name,
		contact.Phone,
		contact.Message,
	)
	return err
}
