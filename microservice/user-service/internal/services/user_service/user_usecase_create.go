package userservice

import (
	"context"
	"errors"
	"user-service/internal/repositories/users"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/apperror"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (u *user) Register(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	// validate user
	{
		userOutput, err := u.userRepositoryReader.FindOne(ctx, users.FindOneInput{
			Email: input.Email,
		})
		if err != nil {
			if !errors.Is(err, databases.ErrNoRowFound) {
				return output, err
			}
		}
		if userOutput.ID != 0 {
			return output, apperror.ErrBadRequest("email is registerd")
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return output, err
	}

	if err = u.tx.DoTxContext(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}, func(ctx context.Context, tx libpgx.RDBMS) (err error) {
		id, err := u.userRepositoryWriter.Create(ctx, users.CreateInput{
			Entity: users.Entity{
				Name:        input.Name,
				Email:       input.Email,
				PhoneNumber: input.PhoneNumber,
				Password:    string(passwordHash),
				TraceParent: observability.ExtractTraceparent(ctx),
			},
			Tx: tx,
		})
		if err != nil {
			return err
		}

		output.ID = id
		return
	}); err != nil {
		return output, err
	}

	return
}
