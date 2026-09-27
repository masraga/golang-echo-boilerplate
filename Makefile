.PHONY: api.yaml validate-api api-docs generate-api generate-backend init generate-wire generate-mocks clean generate-vendor update-api-handler test-update-api-handler install-tools rename_project

api.yaml:
	@swagger-cli bundle app/api/src/main.yaml -o app/api/api.yaml -t yaml

validate-api:
	@swagger-cli validate app/api/api.yaml

api-docs:
	@npx --yes @redocly/cli build-docs app/api/api.yaml -o api-docs.html

external-api.yaml:
	@for dir in app/externalapi/*/; do \
		pkg=$$(basename $$dir); \
		if [ -f "app/externalapi/$$pkg/src/main.yaml" ]; then \
			echo "Bundling $$pkg"; \
			npx swagger-cli bundle app/externalapi/$$pkg/src/main.yaml -o app/externalapi/$$pkg/api.yaml -t yaml; \
			swagger-cli validate app/externalapi/$$pkg/api.yaml; \
		fi \
	done

generate-external-api.yaml:
	@for dir in app/externalapi/*/; do \
		pkg=$$(basename $$dir); \
		if [ -f "app/externalapi/$$pkg/src/main.yaml" ]; then \
			echo "Generating Go client for $$pkg"; \
			mkdir -p ./generated/externalapi/$$pkg; \
			oapi-codegen \
				-generate types,client \
				-package $$pkg \
					-response-type-suffix APIResponse \
				-o ./generated/externalapi/$$pkg/api.gen.go \
				app/externalapi/$$pkg/api.yaml; \
		fi \
	done

generate-external-api-docs:
	@npx --yes @redocly/cli build-docs app/externalapi/$(pkg)/api.yaml -o $(pkg)-docs.html

generate-api:
	@mkdir -p ./generated/api
	@oapi-codegen \
		-generate types,echo-server,spec \
		-package api \
		./app/api/api.yaml > ./generated/api/api.gen.go

update-api-handler: api.yaml validate-api generate-api
	@./bin/update_api_handler.sh

test-update-api-handler:
	@./bin/update_api_handler_test.sh

generate-backend:
	@mkdir -p ./generated/app
	@go build -o ./generated/app/backend ./app/backend

generate-wire:
	wire ./...

LOG_DIR=./internal/log
LOG_FILE=$(LOG_DIR)/error.log
generate-log:
	@mkdir -p ./internal/log
	@touch $(LOG_FILE)

INTERFACE_GO_FILES := $(shell find internal -type f -name "interface.go")
INTERFACE_MOCK_GO_FILES := $(INTERFACE_GO_FILES:%.go=%.mock.gen.go)

# Generate mocks for interfaces
generate-mocks: $(INTERFACE_MOCK_GO_FILES)

$(INTERFACE_MOCK_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	@mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

generate-vendor:
	go mod tidy

init: api.yaml \
	generate-api \
	generate-external-api.yaml \
	generate-wire \
	generate-backend \
	generate-mocks \
	generate-log \
	generate-vendor

clean:
	@rm -rf ./api-docs.html
	@rm -rf ./generated
	@rm -rf ./app/backend/wire_gen.go
	@find . -name "*.mock.gen.go" -type f -delete

migrate_up:
	@migrate -path ./migrations -database $(db) up

migrate_down:
	@migrate -path ./migrations -database $(db) down

migrate_version:
	@migrate -path ./migrations -database $(db) version

migrate_force:
	@migrate -path ./migrations -database $(db) force $(version)

create_migration:
	@mkdir -p ./migrations
	@migrate create -ext sql -dir migrations -format 200601021504 -tz Asia/Jakarta $(filter-out $@,$(MAKECMDGOALS))

ifneq ($(filter create_migration,$(MAKECMDGOALS)),)
MIGRATION_NAME := $(filter-out create_migration,$(MAKECMDGOALS))
ifneq ($(strip $(MIGRATION_NAME)),)
$(MIGRATION_NAME):
	@:
endif
endif

install_tools:
	@bash ./bin/install_dependecies.sh

rename_project:
ifndef name
	$(error name parameter is required. Usage: make rename_project name=github.com/user/new-repo)
endif
	go mod edit -module $(name)
	find . -type f -name "*.go" -exec sed -i '' 's|github.com/masraga/golang-echo-boilerplate|$(name)|g' {} +
	go mod tidy