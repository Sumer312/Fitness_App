include .env

run:
	@./Health-App-Backend
build:
	@templ generate 
	@go build
migrateUp:
	@cd ./sql/schemas && @goose $(DB_TYPE) $(DB_URL) up
	@cd ../.. && @sqlc generate
migrateDown:
	@cd ./sql/schemas && @goose $(DB_TYPE) $(DB_URL) down
	@cd ../..
