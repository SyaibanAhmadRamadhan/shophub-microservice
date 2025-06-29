#!/usr/bin/env bash
# This script is used to create migration the database schema.

# npx openapi-format api/openapi/api.yaml -s api/openapi/openapi-sort.json -f api/openapi/openapi-filter.json -o api/openapi/api.yaml
mkdir -p .gen/api
go tool oapi-codegen --package api -generate types api/openapi/api.yaml > .gen/api/api-types.gen.go
go tool oapi-codegen --package api -generate gin,spec api/openapi/api.yaml > .gen/api/api-server.gen.go