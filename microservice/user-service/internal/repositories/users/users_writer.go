package users

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
)

func (r *user) Create(ctx context.Context, input CreateInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create user, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		`INSERT INTO users (
			name, email, phone_number, password, is_verified, trace_parent
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		input.Entity.Name,
		input.Entity.Email,
		input.Entity.PhoneNumber,
		input.Entity.Password,
		input.Entity.IsVerified,
		input.Entity.TraceParent,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return
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
