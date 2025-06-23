.PHONY: clean \
		generate_api \
		preview_open_api

clean:
	rm -rf dist/* .gen build vendor
	find . -name "*.mock.gen.go" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "wire_gen.go" -type f -delete
	find . -name "*.mock.gen.go" -type f -delete


generate: generate_api

generate_api: api/openapi/api.yaml
	npx openapi-format api/openapi/api.yaml -s api/openapi/openapi-sort.json -f api/openapi/openapi-filter.json -o api/openapi/api.yaml
	mkdir -p .gen/api
	go tool oapi-codegen --package api -generate types $< > .gen/api/api-types.gen.go
	go tool oapi-codegen --package api -generate gin,spec $< > .gen/api/api-server.gen.go

preview_open_api:
	redocly preview-docs api/openapi/api.yaml