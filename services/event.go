package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/kazimanzurrashid/consents-api-go/models"
)

type Event interface {
	Create(
		ctx context.Context,
		request *models.EventCreateRequest) error
}

type PostgresEvent struct {
	db *sql.DB
}

func NewEvent(db *sql.DB) Event {
	return &PostgresEvent{db}
}

func (e *PostgresEvent) Create(
	ctx context.Context,
	request *models.EventCreateRequest) error {
	const query = `INSERT INTO "events"(id, user_id, consent_id, created_at, enabled) VALUES($1, $2, $3, $4, $5)`

	tx, err := e.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})

	if err != nil {
		return err
	}

	for _, consent := range *request.Consents {
		if _, err := tx.ExecContext(
			ctx,
			query,
			generateID(),
			request.User.ID,
			consent.ID,
			time.Now().Format(time.RFC3339),
			consent.Enabled); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
