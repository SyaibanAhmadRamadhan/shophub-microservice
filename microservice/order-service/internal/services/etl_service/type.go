package etlservice

import (
	"order-service/internal/repositories/users"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type EtlUserEntity struct {
	users.Entity
	primitive.DebeziumExtractNewRecordStatePayloadMetadata
}
