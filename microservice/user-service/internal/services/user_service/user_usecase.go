package userservice

import (
	useraddresses "user-service/internal/repositories/user_addresses"
	"user-service/internal/repositories/users"

	libpgx "github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/pgx"
)

type user struct {
	userRepositoryReader        users.RepositoryReader
	userRepositoryWriter        users.RepositoryWriter
	userAddressRepositoryWriter useraddresses.RepositoryWriter
	userAddressRepositoryReader useraddresses.RepositoryReader
	tx                          libpgx.Tx
}

type OptionParams struct {
	UserRepositoryReader        users.RepositoryReader
	UserRepositoryWriter        users.RepositoryWriter
	UserAddressRepositoryWriter useraddresses.RepositoryWriter
	UserAddressRepositoryReader useraddresses.RepositoryReader
	Tx                          libpgx.Tx
}

func New(optionParams OptionParams) *user {
	return &user{
		userRepositoryReader:        optionParams.UserRepositoryReader,
		userRepositoryWriter:        optionParams.UserRepositoryWriter,
		userAddressRepositoryWriter: optionParams.UserAddressRepositoryWriter,
		userAddressRepositoryReader: optionParams.UserAddressRepositoryReader,
		tx:                          optionParams.Tx,
	}
}
