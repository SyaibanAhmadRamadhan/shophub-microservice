package useraddresses

import (
	"context"
	"errors"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

func (r *userAddresses) UpSert(ctx context.Context, input UpSertInput) (id int64, err error) {
	if input.Tx == nil {
		return id, errors.New("failed to create user, transaction database is nil")
	}

	err = input.Tx.QueryRow(ctx,
		`INSERT INTO user_addresses (
			id, user_id, full_address, trace_parent, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			full_address = EXCLUDED.full_address,
			trace_parent = EXCLUDED.trace_parent,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id`,
		input.Entity.ID,
		input.Entity.UserID,
		input.Entity.FullAddress,
		input.Entity.TraceParent,
		input.Entity.CreatedAt,
		input.Entity.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
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

func (r *userAddresses) Delete(ctx context.Context, tx libpgx.RDBMS, id int64) error {
	if tx == nil {
		return errors.New("failed to delete user, transaction database is nil")
	}

	cmdTag, err := tx.Exec(ctx, `DELETE FROM user_addresses WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return databases.ErrNoDeleteRow
	}

	return nil
}
