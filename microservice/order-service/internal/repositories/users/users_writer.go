package users

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

func (r *user) UpSert(ctx context.Context, input UpSertInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create user, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		`INSERT INTO users (
			id, name, email, phone_number, password, is_verified, trace_parent, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			phone_number = EXCLUDED.phone_number,
			password = EXCLUDED.password,
			is_verified = EXCLUDED.is_verified,
			trace_parent = EXCLUDED.trace_parent,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id`,
		input.Entity.ID,
		input.Entity.Name,
		input.Entity.Email,
		input.Entity.PhoneNumber,
		input.Entity.Password,
		input.Entity.IsVerified,
		input.Entity.TraceParent,
		input.Entity.CreatedAt,
		input.Entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *user) Update(ctx context.Context, input UpdateInput) error {
	if input.Tx == nil {
		return errors.New("failed to update user, transaction database is nil")
	}

	cmdTag, err := input.Tx.Exec(ctx,
		`UPDATE users 
		 SET 
			name = $1,
			email = $2,
			phone_number = $3,
			password = $4,
			is_verified = $5,
			trace_parent = $6,
			updated_at = $7
		 WHERE id = $8`,
		input.Entity.Name,
		input.Entity.Email,
		input.Entity.PhoneNumber,
		input.Entity.Password,
		input.Entity.IsVerified,
		input.Entity.TraceParent,
		time.Now().UTC(),
		input.Entity.ID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return databases.ErrNoUpdateRow
	}

	return nil
}
func (r *user) Delete(ctx context.Context, tx libpgx.RDBMS, id int64) error {
	if tx == nil {
		return errors.New("failed to delete user, transaction database is nil")
	}

	cmdTag, err := tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return databases.ErrNoDeleteRow
	}

	return nil
}
