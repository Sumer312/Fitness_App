run:
	@./Fitness
install:
	@go mod tidy && go install github.com/a-h/templ/cmd/templ@v0.2.707
build:
	@templ generate 
	@go build -o Fitness
migrateUp:
	@cd ./sql/schemas && goose $(DB_TYPE) $(DB_URL) up
migrateDown:
	@cd ./sql/schemas && goose $(DB_TYPE) $(DB_URL) down
queries:
	@sqlc generate
