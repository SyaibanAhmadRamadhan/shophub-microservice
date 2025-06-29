package etlservice

import (
	useraddresses "order-service/internal/repositories/user_addresses"
	"order-service/internal/repositories/users"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type EtlUserEntity struct {
	users.Entity
	primitive.DebeziumExtractNewRecordStatePayloadMetadata
}

type EtlUserAddressEntity struct {
	useraddresses.Entity
	primitive.DebeziumExtractNewRecordStatePayloadMetadata
}
