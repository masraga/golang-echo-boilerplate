.PHONY: generate-api generate-backend init generate-wire generate-mocks clean

api.yaml:
	@swagger-cli bundle app/api/src/main.yaml -o app/api/api.yaml -t yaml

validate-api:
	@swagger-cli validate app/api/api.yaml

api-docs:
	@redocly build-docs app/api/api.yaml -o api-docs.html

generate-api:
	@mkdir -p ./generated/api
	@oapi-codegen \
		-generate types,echo-server,spec \
		-package api \
		./app/api/api.yaml > ./generated/api/api.gen.go

generate-backend:
	@mkdir -p ./generated/app
	@go build -o ./generated/app/backend ./app/backend

generate-wire:
	wire ./...

INTERFACE_GO_FILES := $(shell find internal -type f -name "interface.go")
INTERFACE_MOCK_GO_FILES := $(INTERFACE_GO_FILES:%.go=%.mock.gen.go)

# Generate mocks for interfaces
generate-mocks: $(INTERFACE_MOCK_GO_FILES)

$(INTERFACE_MOCK_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	@mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

init: api.yaml generate-api generate-wire generate-backend generate-mocks

clean:
	@rm -rf ./api-docs.html
	@rm -rf ./generated
	@rm -rf ./app/backend/wire_gen.go 
	@find . -name "*.mock.gen.go" -type f -delete

migrate_up:
	@migrate -path ./migrations -database $(db) up

migrate_down:
	@migrate -path ./migrations -database $(db) up

migrate_version:
	@migrate -path ./migrations -database $(db) version

migrate_force:
	@migrate -path ./migrations -database $(db) force $(version)

create_migration:
	@mkdir -p ./migrations
	@migrate create -ext sql -dir migrations -format 200601021504 -tz Asia/Jakarta $(filter-out $@,$(MAKECMDGOALS))
