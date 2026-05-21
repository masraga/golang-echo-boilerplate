.PHONY: generate-api generate-backend init generate-wire generate-mocks

api.yaml:
	@swagger-cli bundle app/api/src/main.yaml -o app/api/api.yaml -t yaml

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

generate-mocks:
	@find . -type f -name "interface.go" | while read f; do \
		dir=$$(dirname "$$f"); \
		base=$$(basename "$$f" .go); \
		mock_file="$$dir/$${base}.mock.gen.go"; \
		echo "Generating mock for $$mock_file"; \
		mockgen -source="$$f" -destination="$$mock_file" -package=$$(basename "$$dir"); \
	done

init: api.yaml generate-api generate-mocks generate-wire generate-backend

clean:
	@rm -rf ./generated
	@rm -rf ./app/backend/wire_gen.go 
	@find . -type f -name "interface.go" | while read f; do \
		dir=$$(dirname "$$f"); \
		base=$$(basename "$$f" .go); \
		mock_file="$$dir/$${base}.mock.gen.go"; \
		rm -f "$$mock_file"; \
	done

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