package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yyek0/stroydom-website/internal/models"
)

type LeadStorage interface {
	Ping(ctx context.Context) error
	Create(ctx context.Context, lead models.Lead) (int, error)
	GetAll(ctx context.Context) ([]models.Lead, error)
	Get(ctx context.Context, id int) (models.Lead, error)
	Delete(ctx context.Context, id int) error
}

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, connString string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		Pool: pool,
	}, nil
}

func (db *PostgresDB) Init(ctx context.Context) error {
	sqlQuery := `
        CREATE TABLE IF NOT EXISTS leads (
            id SERIAL PRIMARY KEY,
            name VARCHAR(200),
            phone VARCHAR(20),
            created_at TIMESTAMP
        );
    `

	_, err := db.Pool.Exec(ctx, sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) Ping(ctx context.Context) error {
	if err := db.Pool.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (db *PostgresDB) Create(ctx context.Context, lead models.Lead) (int, error) {
	sqlQuery := `
		INSERT INTO leads (name, phone, created_at) 
		VALUES ($1, $2, $3) RETURNING id
	`

	var id int
	err := db.Pool.QueryRow(ctx, sqlQuery, lead.Name, lead.Phone, lead.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *PostgresDB) GetAll(ctx context.Context) ([]models.Lead, error) {
	sqlQuery := `
		SELECT id, name, phone, created_at 
		FROM leads
	`
	var leads []models.Lead

	rows, err := db.Pool.Query(ctx, sqlQuery)

	if err != nil {
		return leads, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id        int
			name      string
			phone     string
			createdat time.Time
		)

		if err := rows.Scan(
			&id,
			&name,
			&phone,
			&createdat,
		); err != nil {
			// log it
			return leads, err
		}
		leads = append(leads, models.Lead{
			ID:        id,
			Name:      name,
			Phone:     phone,
			CreatedAt: createdat,
		})
	}

	return leads, nil

}

func (db *PostgresDB) Delete(ctx context.Context, id int) error {
	sqlQuery := `
		DELETE FROM leads WHERE id = $1 
	`

	if _, err := db.Pool.Exec(ctx, sqlQuery, id); err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) Get(ctx context.Context, id int) (models.Lead, error) {
	sqlQuery := `
		SELECT name, phone, created_at 
		FROM leads
		WHERE id = $1
	`

	row := db.Pool.QueryRow(ctx, sqlQuery, id)

	var (
		name      string
		phone     string
		createdAt time.Time
	)

	if err := row.Scan(
		&name,
		&phone,
		&createdAt,
	); err != nil {
		return models.Lead{}, err
	}

	return models.Lead{
		ID:        id,
		Name:      name,
		Phone:     phone,
		CreatedAt: createdAt,
	}, nil
}
