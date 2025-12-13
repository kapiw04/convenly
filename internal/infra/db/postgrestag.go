package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/kapiw04/convenly/internal/domain/event"
)

type PostgresTagRepo struct {
	DB *sql.DB
}

func NewPostgresTagRepo(db *sql.DB) *PostgresTagRepo {
	repo := &PostgresTagRepo{DB: db}
	if err := repo.SeedDefaults(); err != nil {
		panic("failed to seed default tags: " + err.Error())
	}
	return repo
}

func (r *PostgresTagRepo) FindAll() ([]event.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := "SELECT tag_id, name FROM tags"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []event.Tag
	for rows.Next() {
		var t event.Tag
		if err := rows.Scan(&t.TagID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, rows.Err()
}

func (r *PostgresTagRepo) FindByName(name string) (*event.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := "SELECT tag_id, name FROM tags WHERE name = $1"
	row := r.DB.QueryRowContext(ctx, query, name)

	var t event.Tag
	if err := row.Scan(&t.TagID, &t.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTagRepo) CreateIfNotExists(name string) (*event.Tag, error) {
	existing, err := r.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO tags (name) 
		VALUES ($1) 
		RETURNING tag_id, name`

	row := r.DB.QueryRowContext(ctx, query, name)

	var t event.Tag
	err = row.Scan(&t.TagID, &t.Name)
	if err != nil {
		if existing, findErr := r.FindByName(name); findErr == nil && existing != nil {
			return existing, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTagRepo) SeedDefaults() error {
	for _, tagName := range event.DefaultTagNames {
		if _, err := r.CreateIfNotExists(tagName); err != nil {
			return err
		}
	}
	return nil
}
