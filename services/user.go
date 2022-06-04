package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kazimanzurrashid/consents-api-go/models"
)

type User interface {
	Create(
		ctx context.Context,
		request *models.UserCreateRequest) (*models.User, error)

	Delete(ctx context.Context, id string) error

	Detail(ctx context.Context, id string) (*models.User, error)
}

type PostgresUser struct {
	db *sql.DB
}

func NewUser(db *sql.DB) User {
	return &PostgresUser{db}
}

func (u *PostgresUser) Create(
	ctx context.Context,
	request *models.UserCreateRequest) (*models.User, error) {

	const query = `INSERT INTO "users"(id, email) VALUES($1, $2)`
	id := generateID()

	if _, err := u.db.ExecContext(ctx, query, id, request.Email); err != nil {
		return nil, err
	}

	return &models.User{
		ID:       id,
		Email:    request.Email,
		Consents: make([]models.Consent, 0),
	}, nil
}

func (u *PostgresUser) Delete(ctx context.Context, id string) error {
	const eventsQuery = `DELETE FROM "events" WHERE user_id = $1`
	const userQuery = `DELETE FROM "users" WHERE id = $1`

	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})

	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, eventsQuery, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, userQuery, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (u *PostgresUser) Detail(
	ctx context.Context,
	id string) (*models.User, error) {

	const userQuery = `SELECT id, email FROM "users" WHERE id = $1`
	const eventsQuery = `
(SELECT consent_id, enabled
FROM "events"
WHERE user_id = $1
AND consent_id = $%v
ORDER BY created_at
DESC LIMIT 1)`

	userRow := u.db.QueryRowContext(ctx, userQuery, id)

	var user models.User

	if err := userRow.Scan(&user.ID, &user.Email); err != nil {
		return nil, nil
	}

	var consentsQuery string
	values := []interface{}{id}

	for index, consentID := range []string{
		models.ConsentEmail,
		models.ConsentSMS} {
		if index > 0 {
			consentsQuery += " UNION ALL "
		}
		consentsQuery += fmt.Sprintf(eventsQuery, index+2)
		values = append(values, consentID)
	}

	eventRows, err := u.db.QueryContext(ctx, consentsQuery, values...)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = eventRows.Close()
	}()

	consents := make([]models.Consent, 0)

	for eventRows.Next() {
		var consent models.Consent

		if err := eventRows.Scan(&consent.ID, &consent.Enabled); err != nil {
			return nil, err
		}

		consents = append(consents, consent)
	}

	user.Consents = consents

	return &user, nil
}
