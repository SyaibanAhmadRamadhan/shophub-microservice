package useraddresses

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
)

func (r *userAddresses) Create(ctx context.Context, input CreateInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create user, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		"insert into user_addresses (user_id, full_address, trace_parent) VALUES ($1, $2, $3)",
		input.Entity.UserID,
		input.Entity.FullAddress,
		input.Entity.TraceParent,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return
}

func (r *userAddresses) Update(ctx context.Context, input UpdateInput) error {
	if input.Tx == nil {
		return errors.New("failed to update user address, transaction database is nil")
	}

	cmdTag, err := input.Tx.Exec(ctx,
		`UPDATE user_addresses 
		 SET user_id = $1, full_address = $2 , trace_parent = $3, updated_at = $4
		 WHERE id = $5`,
		input.Entity.UserID,
		input.Entity.FullAddress,
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
