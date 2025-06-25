package authservice

import (
	"user-service/internal/repositories/users"
)

type auth struct {
	userRepositoryReader users.RepositoryReader
}

type OptionParams struct {
	UserRepositoryReader users.RepositoryReader
}

func New(optionParams OptionParams) *auth {
	return &auth{
		userRepositoryReader: optionParams.UserRepositoryReader,
	}
}
