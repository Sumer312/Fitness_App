run:
	@./Health-App-Backend
install:
	@go mod vendor && go install github.com/a-h/templ/cmd/templ@latest
build:
	@templ generate 
	@go build
migrateUp:
	@cd ./sql/schemas && goose $(DB_TYPE) $(DB_URL) up
migrateDown:
	@cd ./sql/schemas && goose $(DB_TYPE) $(DB_URL) down
queries:
	@sqlc generate
